// Package env contains enviroment utility functions
package env

import (
	"os"
	"strconv"
)

func GetString(
	key string,
	defaultValue string,
) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func GetNumber(
	key string,
	defaultValue int,
) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	iVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return iVal
}
