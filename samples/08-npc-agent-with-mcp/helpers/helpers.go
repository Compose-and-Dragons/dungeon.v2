package helpers

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func StringToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Cannot convert to int:", err)
		return 0
	}
	return num
}

func StringToFloat(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Cannot convert to float:", err)
		return 0.0
	}
	return num
}

func StringToBool(str string) bool {
	val, err := strconv.ParseBool(str)
	if err != nil {
		fmt.Println("Cannot convert to bool:", err)
		return false
	}
	return val
}
