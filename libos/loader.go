package libos

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/spf13/afero"
	"wetee.app/libos-entry/utils"
)

func PreLoad(chainAddr string, fs afero.Fs) error {
	// 读取配置文件
	// Read config file
	isTee := utils.GetEnv("IN_TEE", "0")
	AppID := utils.GetEnv("APPID", "NONE")
	keyFile := filepath.Join(utils.GetRootDir(), "sid")

	// 读取配置id
	// Read config id
	workerAddr := "https://127.0.0.1:8883"
	if isTee == "1" {
		workerAddr = "https://wetee-worker.worker-system.svc.cluster.local:8883"
	}

	// 读取签名key
	sigKey, err := utils.GetKey(fs, keyFile)
	if err != nil {
		return err
	}

	// 初始化机密注入
	// Initializes the confidential injection
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	param := &LoadParam{
		Address:   sigKey.Address,
		Time:      fmt.Sprint(time.Now().Unix()),
		Signature: "NONE",
	}
	pbt, _ := json.Marshal(param)

	// 签名
	// Sign
	sig, err := signature.Sign([]byte(param.Time), sigKey.URI)
	if err != nil {
		return err
	}
	param.Signature = hex.EncodeToString(sig)

	// 向集群请求机密
	// Request confidential
	bt, err := workerPost(tlsConfig, workerAddr+"/appLoader/"+AppID, string(pbt))
	if err != nil {
		return err
	}

	// 解析机密
	// Parse the secret
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
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}
