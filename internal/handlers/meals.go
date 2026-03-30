package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/akram-fattah/food-delivery/internal/models"
	"strconv"
	"fmt"
	"os"
	"io"
	"time"
)

func CreateMeal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		helper.SendError(w, "فشل قراءة البيانات", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	categoryIDStr := r.FormValue("category_id")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		helper.SendError(w, "السعر غير صحيح", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		helper.SendError(w, "category_id غير صالح", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		helper.SendError(w, "الصورة مطلوبة", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filePath := "uploads/" + filename

	dst, err := os.Create(filePath)
	if err != nil {
		helper.SendError(w, "فشل حفظ الصورة", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		helper.SendError(w, "فشل نسخ الصورة", http.StatusInternalServerError)
		return
	}

	exists, err := database.CategoryExists(context.Background(), categoryID)
	if err != nil || !exists {
		helper.SendError(w, "التصنيف غير موجود", http.StatusBadRequest)
		return
	}

	meal := models.Meal{
		Name:        &name,
		Description: &description,
		Price:       &price,
		ImageURL:    &filePath,
		CategoryID:  &categoryID,
	}

	err = database.CreateMeal(context.Background(), &meal)
	if err != nil {
		helper.SendError(w, "فشل إنشاء الوجبة", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(meal)
}

func GetMeals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	meals, err := database.GetMeals(context.Background())
	if err != nil {
		helper.SendError(w, "فشل جلب الوجبات", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

func GetMealByID(w http.ResponseWriter, r *http.Request) {
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

	meal, err := database.GetMealByID(context.Background(), id)
	if err != nil {
		helper.SendError(w, "الوجبة غير موجودة", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meal)
}

func DeleteMeal(w http.ResponseWriter, r *http.Request) {
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
	err = database.DeleteMeal(context.Background(), id)
	if err != nil {
		helper.SendError(w, "فشل حذف الوجبة", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateMeal(w http.ResponseWriter, r *http.Request) {
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

	var meal models.Meal
	err = json.NewDecoder(r.Body).Decode(&meal)
	if err != nil {
		helper.SendError(w, "بيانات غير صالحة", http.StatusBadRequest)
		return
	}

	if meal.Name == nil &&
		meal.Description == nil &&
		meal.Price == nil &&
		meal.ImageURL == nil &&
		meal.CategoryID == nil &&
		meal.IsAvailable == nil {
		helper.SendError(w, "ما في بيانات للتحديث", http.StatusBadRequest)
		return
	}

	meal.ID = id

	err = database.UpdateMeal(context.Background(), &meal)
	if err != nil {
		helper.SendError(w, "فشل تحديث الوجبة", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meal)
}

func GetMealsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	categoryID, err := helper.GetIDFromURL(r.URL.Path)
	if err != nil {
		helper.SendError(w, "ID غير صالح", http.StatusBadRequest)
		return
	}
	meals, err := database.GetMealsByCategoryID(context.Background(), categoryID)
	if err != nil {
		helper.SendError(w, "فشل جلب الوجبات", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

func SearchMeals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query().Get("q")
	if query == "" {
		helper.SendError(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}
	meals, err := database.SearchMeals(context.Background(), query)
	if err != nil {
		helper.SendError(w, "فشل جلب الوجبات", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}