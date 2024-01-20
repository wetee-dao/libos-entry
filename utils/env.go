package utils

import "os"

func GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}

func GetRootDir() string {
	// 如果是运行在隐私集群中
	// If running in a confidential computing cluster
	wd := os.Getenv("WORKER_ROOT_DIR")
	if len(wd) != 0 {
		return wd
	}
	// 获取当前工作目录
	// Get the current working directory
	wd, err := os.Getwd()
	if err == nil {
		return wd
	}
	panic(err)
}
