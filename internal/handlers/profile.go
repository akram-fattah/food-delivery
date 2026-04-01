package handlers

import (
	"context"
	"encoding/json"
	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"net/http"
	"github.com/akram-fattah/food-delivery/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		helper.SendError(w, "Missing token", http.StatusUnauthorized)
		return
	}

	userID, err := helper.ParseJWT(authHeader)
	if err != nil {
		helper.SendError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	profile, err := database.GetProfile(context.Background(), userID)
	if err != nil {
		helper.SendError(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		w.Header().Set("Allow", "PUT, PATCH")
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		helper.SendError(w, "Missing token", http.StatusUnauthorized)
		return
	}

	userID, err := helper.ParseJWT(authHeader)
	if err != nil {
		helper.SendError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password != nil {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(*req.Password), 14)
		str := string(hashed)
		req.Password = &str
	}

	err = database.UpdateUser(context.Background(), userID, req)
	if err != nil {
		helper.SendError(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "updated successfully",
	})
}
