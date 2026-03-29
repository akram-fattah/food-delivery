package handlers

import (
	"net/http"
	"time"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	u, err := helper.ParseUser(r)
	if err != nil {
		helper.SendError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	if u.Email == "" && u.Username == "" {
		helper.SendError(w, "يجب إدخال البريد الإلكتروني أو اسم المستخدم", http.StatusBadRequest)
		return
	}

	if err := helper.ValidateUser(u); err != nil {
		helper.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Email != "" {
		exists, err := database.IsEmailExists(r.Context(), u.Email)
		if err != nil {
			helper.SendError(w, "خطأ في قاعدة البيانات", http.StatusInternalServerError)
			return
		}
		if exists {
			helper.SendError(w, "البريد الإلكتروني مستخدم", http.StatusConflict)
			return
		}
	}

	if u.Username != "" {
		exists, err := database.IsUsernameExists(r.Context(), u.Username)
		if err != nil {
			helper.SendError(w, "خطأ في قاعدة البيانات", http.StatusInternalServerError)
			return
		}
		if exists {
			helper.SendError(w, "اسم المستخدم مستخدم", http.StatusConflict)
			return
		}
	}

	hashedPassword, err := helper.HashPassword(u.Password)
	if err != nil {
		helper.SendError(w, "خطأ داخلي", http.StatusInternalServerError)
		return
	}
	u.Password = hashedPassword

	code := helper.GenerateCode()
	expires := time.Now().Add(24 * time.Hour)

	err = database.CreateUser(r.Context(), &u, code, expires)

	if err != nil {
		helper.SendError(w, "فشل إنشاء المستخدم", http.StatusInternalServerError)
		return
	}

	body := `
	 <html>
    <head>
        <link href="https://fonts.googleapis.com/css2?family=Tajawal:wght@400;700&display=swap" rel="stylesheet">
    </head>
    <body style="font-family: 'Tajawal', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f8f9fa; padding: 20px; direction: rtl; text-align: right;">
        <div style="max-width: 500px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 15px rgba(0,0,0,0.1); border: 1px solid #e1e1e1;">
            
            <div style="background-color: #d35400; padding: 25px; text-align: center;">
                <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: bold; letter-spacing: 1px; font-family: 'Tajawal', sans-serif;">مطعم الوحدة</h1>
                <p style="color: #ffeb3b; margin: 5px 0 0 0; font-size: 14px; font-family: 'Tajawal', sans-serif;">أصل المذاق اليمني</p>
            </div>
            
            <div style="padding: 40px 30px; text-align: center;">
                <h2 style="color: #2c3e50; margin-bottom: 20px; font-family: 'Tajawal', sans-serif;">مرحباً بك يا عزيزي الزبون!</h2>
                <p style="color: #7f8c8d; font-size: 16px; line-height: 1.6; font-family: 'Tajawal', sans-serif;">
                    سعداء جداً بانضمامك إلينا. أنت على بُعد خطوة واحدة من طلب وجبتك المفضلة. يرجى استخدام الكود أدناه لتفعيل حسابك:
                </p>
                
                   <div style="margin: 30px 0; padding: 20px; background: #fdf2e9; border: 2px dashed #d35400; border-radius: 8px;">
						<span style="font-size: 32px; font-weight: bold; color: #d35400;">%s</span>
					</div>
									
                <p style="color: #95a5a6; font-size: 13px; margin-top: 25px; font-family: 'Tajawal', sans-serif;">
                    * تنتهي صلاحية هذا الكود خلال 24 ساعة. إذا لم تطلب هذا الكود، يمكنك تجاهل هذا الإيميل بأمان.
                </p>
            </div>
            
            <div style="background-color: #f1f1f1; padding: 15px; text-align: center; border-top: 1px solid #eeeeee;">
                <p style="color: #bdc3c7; font-size: 11px; margin: 0; font-family: 'Tajawal', sans-serif;">
                    جميع الحقوق محفوظة © 2026 - فريق تطوير اكرم سوفت
                </p>
            </div>
        </div>
    </body>
    </html>`
	
	go func() {
		_ = helper.SendEmail(u.Email, body, code)
	}()

	helper.SendJSON(w, http.StatusCreated, map[string]string{
		"message": "تم التسجيل بنجاح! تحقق من بريدك الإلكتروني",
	})
}
