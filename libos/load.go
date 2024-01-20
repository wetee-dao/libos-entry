package libos

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"wetee.app/libos-entry/utils"
)

func PreLoad(chainAddr string, fs afero.Fs) error {
	// 读取配置文件
	// Read config file
	isTee := utils.GetEnv("IN_TEE", "0")
	rootID := filepath.Join(utils.GetRootDir(), utils.GetEnv("ROOT_ID", "0"))
	workerAddr := "https://127.0.0.1:443"
	if isTee == "1" {
		workerAddr = "https://wetee-worker.worker-system.svc.cluster.local"
	}
	fmt.Println(workerAddr)
	fmt.Println(rootID)

	// 初始化机密注入
	// Initializes the confidential injection
	bt, err := workerPost(nil, workerAddr, "{}")
	if err != nil {
		return err
	}
	secret := &Secrets{}
	err = json.Unmarshal(bt, secret)
	if err != nil {
		return err
	}

	// 部署机密到运行环境
	// Deploy secrets to the runtime environment
	err = applySecrets(secret, fs)
	if err != nil {
		return err
	}

	return nil
}

func applySecrets(s *Secrets, fs afero.Fs) error {
	// 写入配置文件
	// Write config file
	for path, data := range s.Files {
		bt, _ := base64.StdEncoding.DecodeString(data)
		if err := fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return err
		}
		if err := afero.WriteFile(fs, path, bt, 0o600); err != nil {
			return err
		}
	}

	// 设置环境变量
	// Set environment variables
	for key, value := range s.Env {
		if err := os.Setenv(key, string(value)); err != nil {
			return err
		}
	}

	return nil
}
