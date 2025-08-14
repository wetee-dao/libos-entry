package libos

import (
	"os"

	"github.com/wetee-dao/libos-entry/util"
)

func InitGramineEntry(hostfs util.Fs, isMain bool) (string, error) {
	service := os.Args[0]

	// 初始化配置文件/环境变量
	// Initialize configuration files/environment variables
	_, err := PreLoad(hostfs, isMain)
	if err != nil {
		return "", err
	}

	return service, nil
}
