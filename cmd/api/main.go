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
}

func main() {
	database.ConnectDB()
	HandleFunc()
	http.ListenAndServe(":8000", nil)
}
