package handlers

import (
	"net/http"
	"encoding/json"
	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)


func VerifyResetCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	req := struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Code == "" {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "بيانات غير صالحة. يرجى التأكد من ملء جميع الحقول بشكل صحيح.",
		})
		return
	}

	isValid, err := database.VerifyResetCode(r.Context(), req.Email, req.Code)
	if err != nil || !isValid {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "الكود غير صحيح أو انتهت صلاحيته.",
		})
		return
	}
	helper.SendJSON(w, http.StatusOK, map[string]string{
		"message": "الكود صحيح.",
	})
}