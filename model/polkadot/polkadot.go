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

const cloudAddress = "0x72381a1a0c2858fa134b89b72b054bcb51f80a6a"

func NewContract(url []string, pk *ink.Signer) (*Contract, error) {
	client, err := ink.InitClient(url, false)
	if err != nil {
		return nil, err
	}

	cloud, err := cloud.InitCloudContract(client, cloudAddress)
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
