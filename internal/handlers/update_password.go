package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
)

type UpdatePasswordRequest struct {
	Email           string `json:"email"`
	Code            string `json:"code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Code == "" || req.Password == "" || req.ConfirmPassword == "" {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "تم تحديث كلمة المرور إذا كانت البيانات صحيحة.",
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "كلمتا المرور غير متطابقة.",
		})
		return
	}

	if len(req.Password) < 8 || len(req.Password) > 20 {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "كلمة المرور يجب أن تكون بين 8 و 20 حرفًا.",
		})
		return
	}

	if !helper.IsStrongPassword(req.Password) {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "كلمة المرور يجب أن تحتوي على حرف كبير، حرف صغير، رقم، ورمز خاص.",
		})
		return
	}

	valid, err := database.VerifyResetCode(r.Context(), req.Email, req.Code)
	if err != nil || !valid {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "الكود غير صحيح أو انتهت صلاحيته.",
		})
		return
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "حدث خطأ أثناء معالجة طلبك. يرجى المحاولة مرة أخرى.",
		})
		return
	}

	err = database.UpdateUserPasswordAndClearCode(r.Context(), req.Email, hashed)
	if err != nil {
		helper.SendJSON(w, http.StatusOK, map[string]string{
			"message": "حدث خطأ أثناء معالجة طلبك. يرجى المحاولة مرة أخرى.",
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, map[string]string{
		"message": "تم تحديث كلمة المرور إذا كانت البيانات صحيحة.",
	})
}
