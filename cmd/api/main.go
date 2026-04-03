package main

import (
	"net/http"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/handlers"
	"github.com/akram-fattah/food-delivery/internal/whatsapp"
)

func HandleFunc() {
	http.HandleFunc("/auth/logout", handlers.Logout)
	http.HandleFunc("/auth/register", handlers.Register)
	http.HandleFunc("/auth/reset-password", handlers.ResetPassword)
	http.HandleFunc("/auth/verify-email", handlers.VerifyEmail)
	http.HandleFunc("/auth/login", handlers.Login)
	http.HandleFunc("/auth/update-password", handlers.UpdatePassword)
	http.HandleFunc("/auth/refresh", handlers.RefreshToken)
	http.HandleFunc("/auth/verify-reset-code", handlers.VerifyResetCode)

	http.HandleFunc("/create/categories", handlers.CreateCategory)
	http.HandleFunc("/categories", handlers.GetCategories)
	http.HandleFunc("/categories/", handlers.GetCategoryByID)
	http.HandleFunc("/update/categories/", handlers.UpdateCategory)
	http.HandleFunc("/delete/categories/", handlers.DeleteCategory)
	http.HandleFunc("/create/meal", handlers.CreateMeal)
	http.HandleFunc("/meals/", handlers.GetMealByID)
	http.HandleFunc("/meals/by-category/", handlers.GetMealsByCategory)
	http.HandleFunc("/meals/search", handlers.SearchMeals)

	http.HandleFunc("/meals", handlers.GetMeals)
	http.HandleFunc("/delete/meal/", handlers.DeleteMeal)
	http.HandleFunc("/update/meal/", handlers.UpdateMeal)
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/profile", handlers.GetProfile)
	http.HandleFunc("/profile/update", handlers.UpdateProfile)

	http.HandleFunc("/orders/create", handlers.CreateOrder)
	http.HandleFunc("/orders", handlers.GetOrders)
	http.HandleFunc("/orders/update-status/", handlers.UpdateOrderStatus)

}

func main() {
	database.ConnectDB()
	go whatsapp.StartWhatsAppBot()
	HandleFunc()
	http.ListenAndServe(":8000", nil)
}
