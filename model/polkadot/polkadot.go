package polkadot

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/model/polkadot/cloud"
	chainTypes "github.com/wetee-dao/libos-entry/model/types"
)

// Contract
type Contract struct {
	*ink.ChainClient
	signer *ink.Signer
	cloud  *cloud.Cloud
}

func NewContract(url []string, pk *ink.Signer, params map[string]string) (*Contract, error) {
	client, err := ink.InitClient(url, false)
	if err != nil {
		return nil, err
	}

	cloud, err := cloud.InitCloudContract(client, params["cloud_addr"])
	if err != nil {
		util.LogWithPurple("InitCloudContract", err)
		return nil, err
	}

	return &Contract{
		ChainClient: client,
		signer:      pk,
		cloud:       cloud,
	}, nil
}

func (c *Contract) GetPod(id uint64) (*chainTypes.Pod, error) {
	data, _, err := c.cloud.QueryPod(id, ink.DefaultParamWithOrigin(types.AccountID(c.signer.AccountID())))
	if err != nil {
		return nil, err
	}

	if data.IsNone() {
		return nil, errors.New("pod not found")
	}

	pod := data.V
	images := make([]string, 0)
	for _, container := range pod.F1 {
		images = append(images, string(container.F1.Image))
	}

	return &chainTypes.Pod{Images: images}, nil
}
