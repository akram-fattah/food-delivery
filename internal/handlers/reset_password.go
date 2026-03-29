package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "تم إرسال كود إعادة تعيين كلمة المرور إلى بريدك الإلكتروني.",
		})
		return
	}

	code := helper.GenerateCode()
	expires := time.Now().Add(24 * time.Hour)

	_ = database.SetUserVerificationCode(r.Context(), req.Email, code, expires)

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
                <h2 style="color: #2c3e50; margin-bottom: 20px; font-family: 'Tajawal', sans-serif;">إعادة تعيين كلمة المرور</h2>
                <p style="color: #7f8c8d; font-size: 16px; line-height: 1.6; font-family: 'Tajawal', sans-serif;">
                    لقد طلبت إعادة تعيين كلمة المرور الخاصة بك. يرجى استخدام الكود أدناه لإكمال العملية:
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
		_ = helper.SendEmail(req.Email, body, code)
	}()

	helper.SendJSON(w, http.StatusOK, map[string]string{
		"message": "تم إرسال كود إعادة تعيين كلمة المرور إلى بريدك الإلكتروني إذا كان مسجلاً لدينا.",
	})

}