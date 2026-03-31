package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKeyRefresh = helper.GetJWTKey()

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		helper.SendError(w, "لم يتم العثور على التوكن", http.StatusUnauthorized)
		return
	}

	refreshToken := cookie.Value

	userID, err := database.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		helper.SendError(w, "توكن غير صالح أو منتهي", http.StatusUnauthorized)
		return
	}

	accessExp := time.Now().Add(15 * time.Minute)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     accessExp.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKeyRefresh)
	if err != nil {
		helper.SendError(w, "حصل خطأ ما", http.StatusInternalServerError)
		return
	}

	newRefreshToken := helper.GenerateRefreshToken()

	err = database.UpdateRefreshToken(context.Background(), int64(userID), newRefreshToken)
	if err != nil {
		helper.SendError(w, "فشل تحديث التوكن", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true, 
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessTokenString,
	})
}
