package chain

// import (
// 	chain "github.com/wetee-dao/ink.go"
// 	"github.com/wetee-dao/ink.go/util"
// )

// // Chain
// type Chain struct {
// 	*chain.ChainClient
// 	signer *chain.Signer
// }

// func (c *Chain) Close() {
// 	c.Api.Client.Close()
// }

// func InitChain(url string, pk *chain.Signer) (*Chain, error) {
// 	client, err := chain.ClientInit(url, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	util.LogWithBlue("POD PUBKEY", pk.Address)

// 	return &Chain{
// 		ChainClient: client,
// 		signer:      pk,
// 	}, nil
// }
