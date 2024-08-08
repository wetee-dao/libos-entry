package util

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type TeeTrigger struct {
	Tee       TeeParam
	ClusterId uint64
	Callids   []types.U128
}
