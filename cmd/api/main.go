package main

import (
	"fmt"
	"github.com/akram-fattah/food-delivery/internal/database"
)

func main() {
	database.ConnectDB()
	fmt.Println("Hello, Akram!")
}