package main

import (
	"net/http"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/handlers"
	"github.com/akram-fattah/food-delivery/internal/middleware"
	"github.com/akram-fattah/food-delivery/internal/whatsapp"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/auth/register", handlers.Register)
	mux.HandleFunc("/auth/login", handlers.Login)
	mux.HandleFunc("/auth/verify-email", handlers.VerifyEmail)
	mux.HandleFunc("/auth/reset-password", handlers.ResetPassword)
	mux.HandleFunc("/auth/verify-reset-code", handlers.VerifyResetCode)
	mux.HandleFunc("/auth/refresh", handlers.RefreshToken)

	// Static Files
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Protected Routes (Require Auth)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/auth/logout", handlers.Logout)
	protectedMux.HandleFunc("/auth/update-password", handlers.UpdatePassword)
	protectedMux.HandleFunc("/profile", handlers.GetProfile)
	protectedMux.HandleFunc("/profile/update", handlers.UpdateProfile)
	protectedMux.HandleFunc("/orders/create", handlers.CreateOrder)
	protectedMux.HandleFunc("/orders", handlers.GetOrders)

	// Admin Routes (Require Auth + Admin)
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/create/categories", handlers.CreateCategory)
	adminMux.HandleFunc("/update/categories/", handlers.UpdateCategory)
	adminMux.HandleFunc("/delete/categories/", handlers.DeleteCategory)
	adminMux.HandleFunc("/create/meal", handlers.CreateMeal)
	adminMux.HandleFunc("/delete/meal/", handlers.DeleteMeal)
	adminMux.HandleFunc("/update/meal/", handlers.UpdateMeal)
	adminMux.HandleFunc("/update-role", handlers.UpdateRole)
	adminMux.HandleFunc("/orders/update-status/", handlers.UpdateOrderStatus)

	// Public Data Routes (Optional: can be protected if needed)
	mux.HandleFunc("/categories", handlers.GetCategories)
	mux.HandleFunc("/categories/", handlers.GetCategoryByID)
	mux.HandleFunc("/meals", handlers.GetMeals)
	mux.HandleFunc("/meals/", handlers.GetMealByID)
	mux.HandleFunc("/meals/by-category/", handlers.GetMealsByCategory)
	mux.HandleFunc("/meals/search", handlers.SearchMeals)

	// Apply Middlewares
	mux.Handle("/auth/logout", middleware.AuthMiddleware(http.HandlerFunc(handlers.Logout)))
	mux.Handle("/auth/update-password", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdatePassword)))
	mux.Handle("/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetProfile)))
	mux.Handle("/profile/update", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateProfile)))
	mux.Handle("/orders/create", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateOrder)))
	mux.Handle("/orders", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetOrders)))

	// Apply Admin Middlewares
	mux.Handle("/create/categories", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.CreateCategory))))
	mux.Handle("/update/categories/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.UpdateCategory))))
	mux.Handle("/delete/categories/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.DeleteCategory))))
	mux.Handle("/create/meal", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.CreateMeal))))
	mux.Handle("/delete/meal/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.DeleteMeal))))
	mux.Handle("/update/meal/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.UpdateMeal))))
	mux.Handle("/update-role", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.UpdateRole))))
	mux.Handle("/orders/update-status/", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(handlers.UpdateOrderStatus))))

	// Wrap everything with CORS
	return middleware.CORS(mux)
}

func main() {
	database.ConnectDB()
	go whatsapp.StartWhatsAppBot()
	
	handler := SetupRoutes()
	
	http.ListenAndServe(":8000", handler)
}
