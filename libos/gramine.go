package libos

import (
	"os"

	"github.com/spf13/afero"
	"github.com/wetee-dao/libos-entry/util"
)

func InitGramineEntry(chainAddr string, hostfs afero.Fs, sf util.SecretFunction) (string, error) {
	service := os.Args[0]

	// 初始化配置文件/环境变量
	// Initialize configuration files/environment variables
	// err := PreLoad(chainAddr, hostfs, sf)
	// if err != nil {
	// 	return "", err
	// }

	return service, nil
}
