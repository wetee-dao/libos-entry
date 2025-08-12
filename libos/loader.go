package libos

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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

type InitEnv struct {
	AppID      string
	PodID      uint64
	Files      string
	Encrypts   string
	NameSpace  string
	WorkerAddr string
	ChainAddr  string
}

var DefaultChainUrl string = "ws://192.168.111.105:9944"
var DefaultWorkAddr string = "wetee-worker.worker-system.svc.cluster.local:8883"

func PreLoad(fs util.Fs, isMain bool) error {
	// 读取环境变量
	AppID := util.GetEnv("APPID", "NONE")
	PodID := util.GetEnvU64("PODID", 0)
	Files := util.GetEnv("__FILES__", "{}")
	Encrypts := util.GetEnv("__ENCRYPTS__", "{}")
	NameSpace := util.GetEnv("NAME_SPACE", "")
	WorkerAddr := util.GetEnv("WORKER_ADDR", DefaultWorkAddr)
	ChainAddr := util.GetEnv("CHAIN_ADDR", DefaultChainUrl)

	initEnv := InitEnv{
		AppID:      AppID,
		PodID:      PodID,
		Files:      Files,
		Encrypts:   Encrypts,
		NameSpace:  NameSpace,
		WorkerAddr: WorkerAddr,
		ChainAddr:  ChainAddr,
	}
	return preLoad(fs, isMain, initEnv)
}

func PreLoadFromInitData(fs util.Fs, envs map[string]string, isMain bool) error {
	podId := envs["PODID"]
	i, err := strconv.ParseUint(podId, 10, 64)
	if err != nil {
		return err
	}

	initEnv := InitEnv{
		AppID:      envs["APPID"],
		PodID:      i,
		Files:      envs["__FILES__"],
		Encrypts:   envs["__ENCRYPTS__"],
		NameSpace:  envs["NAME_SPACE"],
		WorkerAddr: envs["WORKER_ADDR"],
		ChainAddr:  envs["CHAIN_ADDR"],
	}

	return preLoad(fs, isMain, initEnv)
}

func preLoad(fs util.Fs, isMain bool, initEnv InitEnv) error {
	inkutil.LogWithGray("WorkerAddr", initEnv.WorkerAddr)
	inkutil.LogWithGray("ChainAddr", initEnv.ChainAddr)

	// 生成 本次部署 Key
	podKey, priv, _, err := util.GenerateKeyPair(rand.Reader)
	if err != nil {
		return errors.New("GenerateKeyPair: " + err.Error())
	}

	wChanel, workerReportBt, err := NewTEEClient(initEnv.WorkerAddr)
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
	fmt.Println(string(workerReportBt))
	err = protoio.ReadMessage(bytes.NewBuffer(workerReportBt), workerReport)
	if err != nil {
		return errors.New("Read Report Message: " + err.Error())
	}

	// 初始化区块链链接
	c, err := chain.InitChain(initEnv.ChainAddr, podKey)
	if err != nil {
		return errors.New("Chain.InitChain: " + err.Error())
	}

	// 验证远程worker的证书
	// Verify the worker tee report
	_, err = fs.VerifyReport(workerReport)
	if err != nil {
		c.Close()
		return errors.New("VerifyReport: " + err.Error())
	}
	c.Close()

	// init env
	files := map[string]string{}
	encrypts := map[string]uint64{}
	idNameMap := map[uint64]string{}
	json.Unmarshal([]byte(initEnv.Files), &files)
	json.Unmarshal([]byte(initEnv.Encrypts), &encrypts)

	ids := []uint64{}
	for name, id := range encrypts {
		ids = append(ids, id)
		idNameMap[id] = name
	}

	// 构建签名证明自己在集群中的身份
	// Build the signature to prove your identity in the cluster
	ns, _ := hex.DecodeString(initEnv.NameSpace)
	podMint := &model.TeeCall{
		Time: time.Now().Unix(),
		Tx: &model.TeeCall_PodStart{
			PodStart: &model.PodStart{
				Id:        uint64(initEnv.PodID),
				AppId:     []byte(initEnv.AppID),
				NameSpace: ns,
				PubKey:    podKey.Public(),
				Indexs:    ids,
			},
		},
	}

	// issue report of TEE CALL
	err = fs.IssueReport(*podKey, podMint)
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
		startTEEServer(fs, podKey, initEnv.ChainAddr)
	} else {
		go startTEEServer(fs, podKey, initEnv.ChainAddr)
	}

	return nil
}
