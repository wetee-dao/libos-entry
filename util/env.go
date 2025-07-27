package util

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

func GetEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		// ... handle error
		panic(err)
	}
	return i
}

func GetEnvU64(key string, defaultValue uint64) uint64 {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	i, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		// ... handle error
		panic(err)
	}
	return i
}

func GetRootDir() string {
	return "/wetee"
}
