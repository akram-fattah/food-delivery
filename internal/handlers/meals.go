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
		Name:        name,
		Description: description,
		Price:       price,
		ImageURL:    filePath,
		CategoryID:  categoryID,
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