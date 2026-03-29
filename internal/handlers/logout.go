package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusInternalServerError)
		return
	}
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.SendError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	err := database.DeleteRefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		helper.SendError(w, "فشل في تسجيل الخروج", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "تم تسجيل الخروج بنجاح"})
}
