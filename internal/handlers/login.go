package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/akram-fattah/food-delivery/internal/models"
)

var jwtKey = helper.GetJWTKey()

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.SendError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	var user models.User
	query := "SELECT id, password, role, is_verified FROM users WHERE email = $1"
	err := database.DB.QueryRow(query, input.Email).Scan(&user.ID, &user.Password, &user.Role, &user.IsVerified)
	if err != nil {
		helper.SendError(w, "البريد الإلكتروني أو كلمة المرور خطأ", http.StatusUnauthorized)
		return
	}

	if !user.IsVerified {
		helper.SendError(w, "حسابك غير مفعل. يرجى تفعيل الحساب عبر البريد الإلكتروني", http.StatusForbidden)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		helper.SendError(w, "البريد الإلكتروني أو كلمة المرور خطأ", http.StatusUnauthorized)
		return
	}

	accessExp := time.Now().Add(15 * time.Minute)
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     accessExp.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		helper.SendError(w, "حصل خطأ ما", http.StatusInternalServerError)
		return
	}

	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		helper.SendError(w, "حصل خطأ ما", http.StatusInternalServerError)
		return
	}
	refreshTokenString := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	err = database.SaveRefreshToken(context.Background(), user.ID, refreshTokenString, refreshExp)
	if err != nil {
		helper.SendError(w, "حصل خطأ ما", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
		"role":         user.Role,
	})
}
