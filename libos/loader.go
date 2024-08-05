package libos

import (
	"crypto/ed25519"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/libos-entry/util"
)

func PreLoad(fs util.Fs) error {
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

	workerReport := util.TeeParam{}
	err = json.Unmarshal(workerReportWrap, &workerReport)
	if err != nil {
		return errors.Wrap(err, "Unmarshal worker report")
	}

	// 验证远程worker的证书
	_, err = VerifyReport(&workerReport, fs)
	if err != nil {
		return errors.Wrap(err, "VerifyReport")
	}

	// 验证worker的版本是否与链上一致

	// 验证worker的singer是否与链上一致

	// 生成 本次部署 Key
	_, deployKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return errors.Wrap(err, "ed25519.GenerateKey")
	}

	singer, err := core.Ed25519PairFromPk(deployKey, 42)
	if err != nil {
		return errors.Wrap(err, "core.Ed25519PairFromPk")
	}

	// 获取本地证书
	// Get local certificate
	report, time, err := fs.IssueReport(&singer, nil)
	if err != nil {
		return errors.Wrap(err, "GetRemoteReport")
	}

	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	param := &util.TeeParam{
		Address: singer.SS58Address(42),
		Time:    time,
		Report:  report,
		Data:    nil,
	}

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
