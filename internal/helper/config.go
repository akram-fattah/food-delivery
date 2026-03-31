package helper

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

func GetJWTKey() []byte {
	_ = godotenv.Load()
	secrit := os.Getenv("JWT_SECRET")
	if secrit == "" {
		log.Fatal("JWT_SECRET is not set in environment variables")
	}
	return []byte(secrit)
}