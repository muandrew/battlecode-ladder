package utils

import (
	"os"
	"fmt"
)

func GetRequiredEnv(key string, fail func()) string {
	envKey := "BCL_" + key
	value := os.Getenv(envKey)
	if value == "" {
		fmt.Printf("Your env variable %s was not configured.\n", envKey)
		fail()
	}
	return value
}