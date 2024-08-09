package chain

import (
	"fmt"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/core"
)

// Chain
type Chain struct {
	client *chain.ChainClient
	signer *core.Signer
}

func (c *Chain) Close() {
	c.client.Api.Client.Close()
}

func InitChain(url string, pk *core.Signer) (*Chain, error) {
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return nil, err
	}

	fmt.Println("Node chain pubkey:", pk.Address)

	return &Chain{
		client: client,
		signer: pk,
	}, nil
}

func (c *Chain) Client() *chain.ChainClient {
	return c.client
}
