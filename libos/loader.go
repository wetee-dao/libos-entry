package libos

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"

	"github.com/wetee-dao/libos-entry/libos/chain"
	"github.com/wetee-dao/libos-entry/util"
)

var DefaultChainUrl string = "ws://wetee-node.worker-addon.svc.cluster.local:9944"

func PreLoad(fs util.Fs, isMain bool) error {
	AppID := util.GetEnv("APPID", "NONE")

	// 获取集群中的worker地址
	workerAddr := util.GetEnv("WORKER_ADDR", "https://wetee-worker.worker-system.svc.cluster.local:8883")
	chainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)

	util.LogWithRed("WorkerAddr", workerAddr)

	// 验证远程worker是否可用
	// Initializes the confidential injection
	wChanel := WorkerChannel{TlsConfig: &tls.Config{InsecureSkipVerify: true}}
	workerReportWrap, err := wChanel.Get(workerAddr + "/report")
	if err != nil {
		return errors.New("GetFromWorker report: " + err.Error())
	}

	// 解析远程worker的证书
	// Parse the worker certificate
	workerReport := util.TeeParam{}
	err = json.Unmarshal(workerReportWrap, &workerReport)
	if err != nil {
		return errors.New("Unmarshal worker report: " + err.Error())
	}

	// 生成 本次部署 Key
	deploySinger, err := util.GenerateKeyPair(rand.Reader)
	if err != nil {
		return errors.New("GenerateKeyPair: " + err.Error())
	}

	// 初始化区块链链接
	c, err := chain.InitChain(chainAddr, deploySinger)
	if err != nil {
		c.Close()
		return errors.New("chain.InitChain: " + err.Error())
	}
	_, err = VerifyWorker(&workerReport, fs, c.ChainClient)
	if err != nil {
		c.Close()
		return errors.New("VerifyReport: " + err.Error())
	}
	c.Close()

	// 获取本地证书
	// Get local certificate
	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	param, err := fs.IssueReport(deploySinger, nil)
	if err != nil {
		return errors.New("GetRemoteReport: " + err.Error())
	}
	pbt, _ := json.Marshal(param)

	// 向集群请求机密
	// Request confidential
	bt, err := wChanel.Post(workerAddr+"/appLaunch/"+AppID, string(pbt))
	if err != nil {
		return errors.New("WorkerPost: " + err.Error())
	}

	// 解析机密
	// Parse the secret
	secret := &util.Secrets{}
	err = json.Unmarshal(bt, secret)
	if err != nil {
		return errors.New("Secrets Unmarshal: " + err.Error())
	}

	// 部署机密到运行环境
	// Deploy secrets to the runtime environment
	err = applySecrets(secret, fs)
	if err != nil {
		return errors.New("applySecrets: " + err.Error())
	}

	if isMain {
		startEntryServer(fs, deploySinger, chainAddr)
	} else {
		go startEntryServer(fs, deploySinger, chainAddr)
	}

	return nil
}

func LocalLoad(fs util.Fs, isMain bool) error {
	// 生成 本次部署 Key
	deploySinger, err := util.GenerateKeyPair(rand.Reader)
	if err != nil {
		return errors.New("GenerateKeyPair: " + err.Error())
	}

	if isMain {
		startEntryServer(fs, deploySinger, "")
	} else {
		go startEntryServer(fs, deploySinger, "")
	}

	return nil
}
