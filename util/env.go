package util

import "os"

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

func GetRootDir() string {
	return "/wetee"
}
