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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.RefreshToken == "" {
		helper.SendError(w, "توكن غير صالح", http.StatusBadRequest)
		return
	}

	userID, err := database.GetRefreshToken(context.Background(), input.RefreshToken)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessTokenString,
	})
}
