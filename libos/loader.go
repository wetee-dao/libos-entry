package libos

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/wetee-dao/libos-entry/util"
)

func PreLoad(chainAddr string, fs util.Fs) error {
	isTee := util.GetEnv("IN_TEE", "0")
	AppID := util.GetEnv("APPID", "NONE")

	// 获取集群中的worker地址
	workerAddr := util.GetEnv("WORKER_ADDR", "https://127.0.0.1:8883")
	if isTee == "1" {
		workerAddr = "https://wetee-worker.worker-system.svc.cluster.local:8883"
	}
	util.LogWithRed("WorkerAddr", workerAddr)

	// 验证远程worker是否可用
	// Initializes the confidential injection
	wChanel := WorkerChannel{TlsConfig: &tls.Config{InsecureSkipVerify: true}}
	workerReportWrap, err := wChanel.Get(workerAddr + "/report")
	if err != nil {
		return errors.Wrap(err, "GetFromWorker report")
	}

	workerReport := map[string]string{}
	err = json.Unmarshal(workerReportWrap, &workerReport)
	if err != nil {
		return errors.Wrap(err, "Unmarshal worker report")
	}

	report, err := hex.DecodeString(workerReport["report"])
	if err != nil {
		return errors.Wrap(err, "Hex decode worker report")
	}

	err = fs.VerifyReport(report, nil, nil)
	if err != nil {
		return errors.Wrap(err, "VerifyReport")
	}

	// 获取本地证书
	// Get local certificate
	certBytes, priv, report, err := GetLocalReport(AppID, fs)
	if err != nil {
		return errors.Wrap(err, "GetRemoteReport")
	}

	// 开启机密服务
	// Start the confidential service
	go startEntryServer(certBytes, priv, report)

	// 设置启动密码
	// TODO password 是用户启动时输入
	fs.SetPassword("123456")

	// 读取配置文件
	// Read config file
	keyFile := filepath.Join(util.GetRootDir(), "sid")

	// 读取签名key
	sigKey, err := util.GetKey(fs, keyFile)
	if err != nil {
		return errors.Wrap(err, "Util.GetKey")
	}

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
	bt, err := wChanel.Post(workerAddr+"/appLoader/"+AppID, string(pbt))
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
