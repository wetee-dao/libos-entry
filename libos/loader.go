package libos

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	inkutil "github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/libos/chain"
	"github.com/wetee-dao/libos-entry/model"
	"github.com/wetee-dao/libos-entry/model/protoio"
	"github.com/wetee-dao/libos-entry/util"
	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/suites"
)

var DefaultChainUrl string = "ws://192.168.111.105:9944"
var DefaultWorkAddr string = "wetee-worker.worker-system.svc.cluster.local:8883"

func PreLoad(fs util.Fs, isMain bool) error {
	// 读取环境变量
	AppID := util.GetEnv("APPID", "NONE")
	PodID := util.GetEnvU64("PODID", 0)
	Files := util.GetEnv("__FILES__", "{}")
	Encrypts := util.GetEnv("__ENCRYPTS__", "{}")
	NameSpace := util.GetEnv("NAME_SPACE", "")
	workerAddr := util.GetEnv("WORKER_ADDR", DefaultWorkAddr)
	chainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)

	inkutil.LogWithGray("WorkerAddr", workerAddr)

	// 生成 本次部署 Key
	deploySinger, priv, _, err := util.GenerateKeyPair(rand.Reader)
	if err != nil {
		return errors.New("GenerateKeyPair: " + err.Error())
	}

	wChanel, workerReportBt, err := NewTEEClient(workerAddr)
	if err != nil {
		return errors.New("NewNewClient: " + err.Error())
	}

	go wChanel.Start()

	// 以 ed25519 密钥对生成TLS证书
	// _, _, serverCertDER, err := util.Ed25519Cert(
	// 	"localhost",
	// 	[]net.IP{net.ParseIP("127.0.0.1")},
	// 	[]string{"localhost"},
	// 	priv,
	// 	pub,
	// )
	// if err != nil {
	// 	return errors.New("Ed25519Cert: " + err.Error())
	// }

	// 获取本地证书
	// Get local certificate
	// serverCert := tls.Certificate{Certificate: [][]byte{serverCertDER}, PrivateKey: priv}
	// tlsConfig := &tls.Config{
	// 	Certificates: []tls.Certificate{serverCert},
	// 	MinVersion:   tls.VersionTLS13,

	// 	// skip client verification
	// 	InsecureSkipVerify: true,
	// 	ClientAuth:         tls.RequireAnyClientCert,
	// }

	// 验证远程worker是否可用
	// Initializes the confidential injection
	// workerReportBt, err := wChanel.Get("/report")
	// if err != nil {
	// 	return errors.New("GetFromWorker report: " + err.Error())
	// }

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

	// 验证远程worker的证书
	// Verify the worker tee report
	_, err = fs.VerifyReport(workerReport)
	if err != nil {
		c.Close()
		return errors.New("VerifyReport: " + err.Error())
	}
	c.Close()

	files := map[string]string{}
	encrypts := map[string]uint64{}
	idNameMap := map[uint64]string{}
	json.Unmarshal([]byte(Files), &files)
	json.Unmarshal([]byte(Encrypts), &encrypts)

	ids := []uint64{}
	for name, id := range encrypts {
		ids = append(ids, id)
		idNameMap[id] = name
	}

	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	ns, _ := hex.DecodeString(NameSpace)
	podMint := &model.TeeCall{
		Time: time.Now().Unix(),
		Tx: &model.TeeCall_PodStart{
			PodStart: &model.PodStart{
				Id:        uint64(PodID),
				AppId:     []byte(AppID),
				NameSpace: ns,
				PubKey:    deploySinger.Public(),
				Indexs:    ids,
			},
		},
	}

	// issue report of TEE CALL
	err = fs.IssueReport(*deploySinger, podMint)
	if err != nil {
		return errors.New("GetRemoteReport: " + err.Error())
	}

	// encode msg
	buf := new(bytes.Buffer)
	err = types.WriteMessage(podMint, buf)
	if err != nil {
		return errors.New("WriteMessage: " + err.Error())
	}

	// 向集群请求机密
	// Request confidential
	bt, err := wChanel.Invoke("/launch", buf.Bytes())
	if err != nil {
		return errors.New("Launch: " + err.Error())
	}

	// 解析机密数据
	// Parse the secret data
	secrets := new(model.DecryptResp)
	err = protoio.ReadMessage(bytes.NewBuffer(bt), secrets)
	if err != nil {
		return errors.New("ReadMessage: " + err.Error())
	}

	// 解析机密
	// Parse the secret
	secretEnv := &util.Secrets{
		Envs:  map[string]string{},
		Files: map[string][]byte{},
	}

	// 将机密数据转换为环境变量
	// Convert secret data to environment variables
	suite := suites.MustFind("Ed25519")
	dkgPubKey := suite.Point()
	dkgPubKey.UnmarshalBinary(secrets.DkgKey)
	for index, secret := range secrets.Lists {
		rawXncCmt := secret.XncCmt
		xncCmt := suite.Point()
		xncCmt.UnmarshalBinary(rawXncCmt)
		encScrts := make([]kyber.Point, len(secret.EncScrt))
		for i, rawEncScrt := range secret.EncScrt {
			encScrt := suite.Point()
			encScrt.UnmarshalBinary(rawEncScrt)
			encScrts[i] = encScrt
		}
		data, err := util.DecryptSecret(suite, encScrts, dkgPubKey, xncCmt, util.Ed25519Scalar(suite, priv))
		if err != nil {
			return errors.New("DecryptSecret: " + err.Error())
		}
		secretEnv.Envs[idNameMap[index]] = string(data)
	}

	// 保存文件
	// Save the file
	for key, value := range files {
		bt, _ := hex.DecodeString(value)
		secretEnv.Files[key] = bt
	}

	// 部署机密到运行环境
	// Deploy secrets to the runtime environment
	err = applySecrets(secretEnv, fs)
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
