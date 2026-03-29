package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)

type verifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusInternalServerError)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var input verifyEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helper.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Code == "" {
		helper.SendError(w, "Missing fields", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET is_verified = true,
		    verification_code = NULL,
		    verification_expires = NULL
		WHERE email = $1
		  AND verification_code = $2
		  AND verification_expires > NOW()
		  AND is_verified = false
	`

	result, err := database.DB.ExecContext(ctx, query, input.Email, input.Code)
	if err != nil {
		helper.SendError(w, "Server error", http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		helper.SendError(w, "Server error", http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		helper.SendError(w, "كود غير صحيح أو منتهي الصلاحية أو الحساب مفعل مسبقاً", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "تم تفعيل الحساب بنجاح",
	})
}
