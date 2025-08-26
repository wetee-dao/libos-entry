package model

import (
	"github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
	"github.com/wetee-dao/libos-entry/model/polkadot"
	chainTypes "github.com/wetee-dao/libos-entry/model/types"
)

type Chain interface {
	GetPod(id uint64) (*chainTypes.Pod, error)
}

func ConnectChain(url []string) (Chain, error) {
	pk, err := ink.Sr25519PairFromSecret("//Alice", 42)
	if err != nil {
		util.LogWithPurple("Sr25519PairFromSecret", err)
		panic(err)
	}

	c, err := polkadot.NewContract(url, &pk)
	if err != nil {
		return nil, err
	}
	return c, nil
}
