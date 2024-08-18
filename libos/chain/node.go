package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/weteedsecret"
	"github.com/wetee-dao/go-sdk/pallet/weteeworker"

	"github.com/wetee-dao/libos-entry/util"
)

// RegisterNode register node
// 注册节点
func (c *Chain) RegisterNode(signer *core.Signer, pubkey []byte) error {
	var bt [32]byte
	copy(bt[:], pubkey)

	call := weteedsecret.MakeRegisterNodeCall(bt)
	return c.client.SignAndSubmit(signer, call, true)
}

// 获取节点列表
// GetNodeList get node list
func (c *Chain) GetNodeList() ([][]byte, error) {
	ret, err := c.client.QueryMapAll("WeTEEDsecret", "Nodes")
	if err != nil {
		return nil, err
	}

	nodes := make([][]byte, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := []byte{}
			if err := codec.Decode(change.StorageData, &n); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

// 获取全网当前程序的代码版本
// Get CodeMrenclave
func (c *Chain) GetCodeMrenclave() ([]byte, error) {
	return weteedsecret.GetCodeMrenclaveLatest(c.client.Api.RPC.State)
}

// 获取全网当前程序的签名人
// Get CodeMrsigner
func (c *Chain) GetCodeMrsigner() ([]byte, error) {
	return weteedsecret.GetCodeMrsignerLatest(c.client.Api.RPC.State)
}

// 查询worker列表
// Get WorkerList
func (c *Chain) GetWorkerList() ([]*types.K8sCluster, error) {
	ret, err := c.client.QueryMapAll("WeTEEWorker", "K8sClusters")
	if err != nil {
		return nil, err
	}

	// 获取节点列表
	nodes := make([]*types.K8sCluster, 0)
	for _, elem := range ret {
		for _, change := range elem.Changes {
			n := &types.K8sCluster{}
			if err := codec.Decode(change.StorageData, n); err != nil {
				util.LogWithRed("codec.Decode", err)
				continue
			}
			nodes = append(nodes, n)
		}
	}

	return nodes, nil
}

// 获取worker的BootPeers
// Get BootPeers
func (c *Chain) GetBootPeers() ([]types.P2PAddr, error) {
	return weteeworker.GetBootPeersLatest(c.client.Api.RPC.State)
}
