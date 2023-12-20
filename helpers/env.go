package helpers

import (
	"os"
	"strconv"
)

func GetEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return intValue
	}

	return defaultValue
}

func GetEnvUint(key string, defaultValue uint) uint {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return uint(intValue)
	}

	return defaultValue
}
