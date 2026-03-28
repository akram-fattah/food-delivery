package helper

import (
	"os"

	"github.com/joho/godotenv"
)

func GetJWTKey() []byte {
	_ = godotenv.Load()
	return []byte(os.Getenv("JWT_SECRET"))
}