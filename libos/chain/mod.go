package chain

import (
	"crypto/ed25519"
	"fmt"

	chain "github.com/wetee-dao/go-sdk"
	"github.com/wetee-dao/go-sdk/core"
)

// ChainClient
var ChainClient *Chain

// Chain
type Chain struct {
	client *chain.ChainClient
	signer *core.Signer
}

func InitChain(url string, pk *ed25519.PrivateKey) error {
	client, err := chain.ClientInit(url, true)
	if err != nil {
		return err
	}

	p, err := core.Ed25519PairFromPk(*pk, 42)
	if err != nil {
		return err
	}
	fmt.Println("Node chain pubkey:", p.Address)

	ChainClient = &Chain{
		client: client,
		signer: &p,
	}
	return nil
}
