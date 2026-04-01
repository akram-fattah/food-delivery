package handlers 


import (	
	"net/http"
	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"encoding/json"
	"context"
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