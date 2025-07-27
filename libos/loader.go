package libos

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/libos/chain"
	"github.com/wetee-dao/libos-entry/model"
	"github.com/wetee-dao/libos-entry/model/protoio"
	"github.com/wetee-dao/libos-entry/util"
)

var DefaultChainUrl string = "ws://192.168.111.105:9944"
var DefaultWorkAddr string = "https://wetee-worker.worker-system.svc.cluster.local:8883"

func PreLoad(fs util.Fs, isMain bool) error {
	// 读取环境变量
	AppID := util.GetEnv("APPID", "NONE")
	PodID := util.GetEnvU64("PODID", 0)
	workerAddr := util.GetEnv("WORKER_ADDR", DefaultWorkAddr)
	chainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)

	inkutil.LogWithGray("WorkerAddr", workerAddr)

	// 生成 本次部署 Key
	deploySinger, priv, pub, err := util.GenerateKeyPair(rand.Reader)
	if err != nil {
		return errors.New("GenerateKeyPair: " + err.Error())
	}

	// 以 ed25519 密钥对生成TLS证书
	_, _, serverCertDER, err := util.Ed25519Cert(
		"localhost",
		[]net.IP{net.ParseIP("127.0.0.1")},
		[]string{"localhost"},
		priv,
		pub,
	)
	if err != nil {
		return errors.New("Ed25519Cert: " + err.Error())
	}

	// Golang tls.config
	serverCert := tls.Certificate{Certificate: [][]byte{serverCertDER}, PrivateKey: priv}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS13,

		// skip client verification
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAnyClientCert,
	}

	// 验证远程worker是否可用
	// Initializes the confidential injection
	wChanel := WorkerChannel{TlsConfig: tlsConfig}
	workerReportBt, err := wChanel.Get(workerAddr + "/report")
	if err != nil {
		return errors.New("GetFromWorker report: " + err.Error())
	}

	// 解析远程worker的证书
	// Parse the worker certificate
	workerReport := new(model.TeeCall)
	err = protoio.ReadMessage(bytes.NewBuffer(workerReportBt), workerReport)
	if err != nil {
		return errors.New("Read Report Message: " + err.Error())
	}

	// 初始化区块链链接
	c, err := chain.InitChain(chainAddr, deploySinger)
	if err != nil {
		c.Close()
		return errors.New("chain.InitChain: " + err.Error())
	}
	_, err = fs.VerifyReport(workerReport)
	if err != nil {
		c.Close()
		return errors.New("VerifyReport: " + err.Error())
	}
	c.Close()

	// 获取本地证书
	// Get local certificate
	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	podMint := &model.TeeCall{
		Time: time.Now().Unix(),
		Tx: &model.TeeCall_PodStart{
			PodStart: &model.PodStart{
				Id: uint64(PodID),
			},
		},
	}
	err = fs.IssueReport(*deploySinger, podMint)
	if err != nil {
		return errors.New("GetRemoteReport: " + err.Error())
	}
	buf := new(bytes.Buffer)
	err = types.WriteMessage(podMint, buf)
	if err != nil {
		return errors.New("WriteMessage: " + err.Error())
	}

	// 向集群请求机密
	// Request confidential
	bt, err := wChanel.PostBt(workerAddr+"/appLaunch/"+AppID, buf.Bytes())
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
		startTEEServer(fs, deploySinger, chainAddr)
	} else {
		go startTEEServer(fs, deploySinger, chainAddr)
	}

	return nil
}
