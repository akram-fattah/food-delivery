package whatsapp

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	_ "github.com/lib/pq"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"github.com/akram-fattah/food-delivery/internal/database"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var Client *whatsmeow.Client

func StartWhatsAppBot() {
	dbLog := waLog.Stdout("Database", "INFO", true)
	clientLog := waLog.Stdout("Client", "INFO", true)

	ctx := context.Background()

	container := sqlstore.NewWithDB(database.DB, "postgres", dbLog)

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	Client = whatsmeow.NewClient(deviceStore, clientLog)

	Client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.Connected:
			fmt.Println("✅ Connected to WhatsApp")

		case *events.Disconnected:
			fmt.Println("❌ Disconnected")
		}
	})

	if Client.Store.ID == nil {
		qrChan, _ := Client.GetQRChannel(ctx)

		err = Client.Connect()
		if err != nil {
			log.Fatal(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("📱 Scan this QR:")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Event:", evt.Event)
			}
		}
	} else {
		err = Client.Connect()
		if err != nil {
			log.Fatal(err)
		}
	}

	select {}
}

func SendMessage(phone string, message string) error {
	if Client == nil {
		return fmt.Errorf("whatsapp client not initialized")
	}

	// تنظيف رقم الهاتف
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.TrimSpace(phone)

	if phone == "" {
		return fmt.Errorf("phone number is empty")
	}

	if message == "" {
		return fmt.Errorf("message is empty")
	}

	// إنشاء JID بشكل صحيح
	jid := types.NewJID(phone, types.DefaultUserServer)

	// إضافة timeout حتى لا يتعلق الطلب
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// إرسال الرسالة
	_, err := Client.SendMessage(ctx, jid, &waProto.Message{
		Conversation: &message,
	})

	if err != nil {
		return fmt.Errorf("failed to send whatsapp message to %s: %w", phone, err)
	}

	return nil
}