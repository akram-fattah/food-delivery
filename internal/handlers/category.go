package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/akram-fattah/food-delivery/internal/models"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		helper.SendError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}
	if cat.CategoryName == "" {
		helper.SendError(w, "اسم التصنيف مطلوب", http.StatusBadRequest)
		return
	}
	id, err := database.CreateCategory(context.Background(), cat.CategoryName)
	if err != nil {
		helper.SendError(w, "الصنف موجود بالفعل", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           id,
		"category_name": cat.CategoryName,
		"created_at":   time.Now(),
	})
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	categories, err := database.GetCategories(context.Background())
	if err != nil {
		helper.SendError(w, "خطأ في جلب التصنيفات", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := helper.GetIDFromURL(r.URL.Path)
	if err != nil {
		helper.SendError(w, "ID غير صالح", http.StatusBadRequest)
		return
	}
	
	category, err := database.GetCategoryByID(context.Background(), id)
	if err != nil {
		helper.SendError(w, "التصنيف غير موجود", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}


func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := helper.GetIDFromURL(r.URL.Path)
	if err != nil {
		helper.SendError(w, "ID غير صالح", http.StatusBadRequest)
		return
	}

	var body struct {
		CategoryName string `json:"category_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helper.SendError(w, "بيانات غير صحيحة", http.StatusBadRequest)
		return
	}

	if body.CategoryName == "" {
		helper.SendError(w, "category_name مطلوب", http.StatusBadRequest)
		return
	}

	updated, err := database.UpdateCategory(context.Background(), id, body.CategoryName)
	if err != nil {
		helper.SendError(w, "التصنيف غير موجود أو حدث خطأ", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := helper.GetIDFromURL(r.URL.Path)
	if err != nil {
		helper.SendError(w, "ID غير صالح", http.StatusBadRequest)
		return
	}

	err = database.DeleteCategory(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			helper.SendError(w, "التصنيف غير موجود", http.StatusNotFound)
			return
		}

		helper.SendError(w, "حدث خطأ أثناء الحذف", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}