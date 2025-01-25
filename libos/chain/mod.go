package chain

import (
	"fmt"

	chain "github.com/wetee-dao/go-sdk"
)

// Chain
type Chain struct {
	*chain.ChainClient
	signer *chain.Signer
}

func (c *Chain) Close() {
	c.Api.Client.Close()
}

func InitChain(url string, pk *chain.Signer) (*Chain, error) {
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return nil, err
	}

	fmt.Println("Node chain pubkey:", pk.Address)

	return &Chain{
		ChainClient: client,
		signer:      pk,
	}, nil
}
