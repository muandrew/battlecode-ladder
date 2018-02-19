package utils

import (
	"github.com/joho/godotenv"
	"log"
)

func InitMainEnv() {
	err := godotenv.Load("bcl-env.sh")
	Initialize("BCL_")
	if err != nil {
		log.Fatalf("Error loading .env file; err: %q", err)
	}
}
