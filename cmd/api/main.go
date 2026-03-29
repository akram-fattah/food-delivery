package main

import (
	"net/http"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/handlers"
)

func HandleFunc() {
	http.HandleFunc("/auth/logout", handlers.Logout)
	http.HandleFunc("/auth/register", handlers.Register)
	http.HandleFunc("/auth/reset-password", handlers.ResetPassword)
	http.HandleFunc("/auth/verify-email", handlers.VerifyEmail)
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/update-password", handlers.UpdatePassword)
	http.HandleFunc("/auth/refresh", handlers.RefreshToken)

	http.HandleFunc("/create/categories", handlers.CreateCategory)
	http.HandleFunc("/categories", handlers.GetCategories)
	http.HandleFunc("/categories/", handlers.GetCategoryByID)
	http.HandleFunc("/update/categories/", handlers.UpdateCategory)
	http.HandleFunc("/delete/categories/", handlers.DeleteCategory)
	http.HandleFunc("/meals", handlers.CreateMeal)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

}

func main() {
	database.ConnectDB()
	HandleFunc()
	http.ListenAndServe(":8000", nil)
}
