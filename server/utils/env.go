package utils

import (
	"os"
	"fmt"
)

var prefix = ""
var isDev bool

func Initialize(appPrefix string) {
	prefix = appPrefix
	isDev = getEnv("ENV") == "DEV"
}

func getEnvKey(key string) string {
	return prefix + key
}

func getEnv(key string) string {
	return os.Getenv(getEnvKey(key))
}

func GetRequiredEnv(key string, fail func()) string {
	envKey := getEnvKey(key)
	value := os.Getenv(envKey)
	if value == "" {
		fmt.Printf("Your env variable %s was not configured.\n", envKey)
		fail()
	}
	return value
}

func IsDev() bool {
	return  isDev
}