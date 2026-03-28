package main

import (
	"net/http"

	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/handlers"
)

func HandleFunc() {
	http.HandleFunc("/auth/register", handlers.Register)
	http.HandleFunc("/auth/verify-email", handlers.VerifyEmail)
}

func main() {
	database.ConnectDB()
	HandleFunc()
	http.ListenAndServe(":8000", nil)
}
