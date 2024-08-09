package chain

import (
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/utility"
	"github.com/wetee-dao/go-sdk/pallet/weteebridge"
)

// list tee calls
func (m *Chain) ListTeeCalls(cid uint64, callId []types.U128) ([]*gtypes.TEECall, []types.StorageKey, error) {
	var pallet, method = "WeTEEBridge", "TEECalls"
	calls := make([]interface{}, 0, len(callId))
	for _, id := range callId {
		calls = append(calls, id)
	}
	set, err := m.client.QueryDoubleMapKeys(pallet, method, cid, calls, nil)
	if err != nil {
		return nil, nil, err
	}

	var list []*gtypes.TEECall = make([]*gtypes.TEECall, 0, len(set))
	var keys []types.StorageKey = make([]types.StorageKey, 0, len(set))
	for _, elem := range set {
		for _, change := range elem.Changes {
			var d gtypes.TEECall

			if err := codec.Decode(change.StorageData, &d); err != nil {
				fmt.Println(err)
				continue
			}
			keys = append(keys, change.StorageKey)
			list = append(list, &d)
		}
	}

	return list, keys, nil
}

func (m *Chain) GetMetaApi(w gtypes.WorkId) (gtypes.ApiMeta, error) {
	api, ok, err := weteebridge.GetApiMetasLatest(m.client.Api.RPC.State, w)
	if err != nil {
		return gtypes.ApiMeta{}, err
	}
	if !ok {
		return gtypes.ApiMeta{}, fmt.Errorf("not found")
	}
	return api, nil
}

func (m *Chain) TeeCallback(cid uint64, callId []types.U128, callbacks []TeeCallBack) error {
	calls := make([]gtypes.RuntimeCall, 0, len(callbacks))
	for i, cb := range callbacks {
		var err []byte
		var isErr bool
		if cb.Err != nil {
			err = []byte(cb.Err.Error())
			isErr = true
		}
		call := weteebridge.MakeInkCallbackCall(cid, callId[i], cb.Resp, types.NewU128(*big.NewInt(0)), gtypes.OptionTByteSlice{
			IsSome:       isErr,
			IsNone:       !isErr,
			AsSomeField0: err,
		})
		calls = append(calls, call)
	}

	call := utility.MakeBatchCall(calls)
	return m.client.SignAndSubmit(m.signer, call, true)
}

type TeeCallBack struct {
	Err  error
	Resp []byte
}
