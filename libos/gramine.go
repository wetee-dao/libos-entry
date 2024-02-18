package libos

import (
	"os"

	"github.com/wetee-dao/libos-entry/util"
)

func InitGramineEntry(chainAddr string, hostfs util.Fs) (string, error) {
	service := os.Args[0]

	// 初始化配置文件/环境变量
	// Initialize configuration files/environment variables
	// err := PreLoad(chainAddr, hostfs, sf)
	// if err != nil {
	// 	return "", err
	// }

	return service, nil
}
