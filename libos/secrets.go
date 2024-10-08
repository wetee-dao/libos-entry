package libos

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/wetee-dao/libos-entry/util"
)

func applySecrets(s *util.Secrets, fs util.Fs) error {
	const keyPrePath = "/dev/attestation/keys/"
	// 先写入其他的加密文件需要的解密钥匙
	// Write encrypted key file for other
	for keyPath, data := range s.Files {
		if strings.HasPrefix(keyPath, keyPrePath) {
			if err := fs.WriteFile(keyPath, data, 0); err != nil {
				return err
			}
			delete(s.Files, keyPath)
		}
	}

	// 写入配置文件
	// Write config file
	for path, data := range s.Files {
		if err := fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return err
		}
		if err := fs.WriteFile(path, data, 0o600); err != nil {
			return err
		}
	}

	// 设置环境变量
	// Set environment variables
	for key, value := range s.Envs {
		if key == "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}
