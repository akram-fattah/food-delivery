package handlers

import (
	"net/http"
	"time"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)

func Register(w http.ResponseWriter, r *http.Request) {

	u, err := helper.ParseUser(r)
	if err != nil {
		helper.SendError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	if err := helper.ValidateUser(u); err != nil {
		helper.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := database.IsEmailExists(r.Context(), u.Email)
	if err != nil {
		helper.SendError(w, "خطأ في قاعدة البيانات", http.StatusInternalServerError)
		return
	}

	if exists {
		helper.SendError(w, "البريد الإلكتروني مستخدم", http.StatusConflict)
		return
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

	go func() {
		_ = helper.SendEmail(u.Email, code)
	}()

	helper.SendJSON(w, http.StatusCreated, map[string]string{
		"message": "تم التسجيل بنجاح! تحقق من بريدك الإلكتروني",
	})
}
