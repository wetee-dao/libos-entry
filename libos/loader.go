package libos

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wetee-dao/libos-entry/util"
)

func PreLoad(chainAddr string, fs util.Fs) error {
	isTee := util.GetEnv("IN_TEE", "0")
	AppID := util.GetEnv("APPID", "NONE")
	AppID = "AAIAAAQAAgABAAAAEhX%2F%2F%2F%2BADgAAAAAAAAAAAAsAAAAAAAD%2FAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAPAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAmsEh%2Fd3L1MFaioQfK5EngnowPzcScScPjQtmv5kiXKl3Nf0Tt9LYBC4ou2g%3D"

	certBytes, priv, report, err := GetRemoteReport(AppID, fs)
	if err != nil {
		return errors.Wrap(err, "GetRemoteReport")
	}

	// 开启机密服务
	// Start the confidential service
	go startEntryServer(certBytes, priv, report)

	// 设置启动密码
	fs.SetPassword("123456")

	// 读取配置文件
	// Read config file
	keyFile := filepath.Join(util.GetRootDir(), "sid")

	// 读取配置id
	workerAddr := "https://127.0.0.1:8883"
	if isTee == "1" {
		workerAddr = "https://wetee-worker.worker-system.svc.cluster.local:8883"
	}

	// 读取签名key
	sigKey, err := util.GetKey(fs, keyFile)
	if err != nil {
		return errors.Wrap(err, "Util.GetKey")
	}

	// 初始化机密注入
	// Initializes the confidential injection
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	param := &util.LoadParam{
		Address:   sigKey.SS58Address(42),
		Time:      fmt.Sprint(time.Now().Unix()),
		Cert:      certBytes,
		Report:    report,
		Signature: "",
	}

	// 签名
	// Sign
	sig, err := sigKey.Sign([]byte(param.Time))
	if err != nil {
		return errors.Wrap(err, "SigKey")
	}
	param.Signature = hex.EncodeToString(sig)
	pbt, _ := json.Marshal(param)

	// 向集群请求机密
	// Request confidential
	bt, err := PostToWorker(tlsConfig, workerAddr+"/appLoader/"+AppID, string(pbt))
	if err != nil {
		return errors.Wrap(err, "WorkerPost")
	}

	// 解析机密
	// Parse the secret
	secret := &util.Secrets{}
	err = json.Unmarshal(bt, secret)
	if err != nil {
		return errors.Wrap(err, "Secrets Unmarshal")
	}

	fmt.Println("Secrets: ", secret)

	// 部署机密到运行环境
	// Deploy secrets to the runtime environment
	err = applySecrets(secret, fs)
	if err != nil {
		return errors.Wrap(err, "applySecrets")
	}

	return nil
}

func applySecrets(s *util.Secrets, fs util.Fs) error {
	const atKeyBasePath = "/dev/attestation/keys/"
	// 先写入其他的加密文件需要的解密钥匙
	// Write encrypted file keys
	for keyPath, data := range s.Files {
		if strings.HasPrefix(keyPath, atKeyBasePath) {
			bt, _ := base64.StdEncoding.DecodeString(data)
			if err := fs.WriteFile(keyPath, bt, 0); err != nil {
				return err
			}
			delete(s.Files, keyPath)
		}
	}

	// 写入配置文件
	// Write config file
	for path, data := range s.Files {
		bt, _ := base64.StdEncoding.DecodeString(data)
		if err := fs.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return err
		}
		if err := fs.WriteFile(path, bt, 0o600); err != nil {
			return err
		}
	}

	// 设置环境变量
	// Set environment variables
	for key, value := range s.Env {
		if key == "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return nil
}
