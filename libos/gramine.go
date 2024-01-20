package libos

import (
	"os"

	"github.com/spf13/afero"
)

func InitGramineEntry(hostfs afero.Fs, chainAddr string) (string, error) {
	service := os.Args[0]

	// 初始化配置文件/环境变量
	// Initialize configuration files/environment variables
	err := PreLoad(chainAddr, hostfs)
	if err != nil {
		return "", err
	}

	return service, nil
}
