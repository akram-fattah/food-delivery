package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"os"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/joho/godotenv"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		helper.SendError(w, "لم يتم العثور على الكوكي", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	err = database.DeleteRefreshToken(r.Context(), refreshToken)
	if err != nil {
		helper.SendError(w, "فشل في تسجيل الخروج", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "تم تسجيل الخروج بنجاح",
	})
}