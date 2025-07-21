package chain

// import (
// 	"fmt"
// 	"math/big"

// 	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
// 	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
// 	"github.com/wetee-dao/go-sdk/pallet/bridge"
// 	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
// 	"github.com/wetee-dao/go-sdk/pallet/utility"
// 	"github.com/wetee-dao/libos-entry/util"
// )

// // list tee calls
// func (m *Chain) ListTeeCalls(cid uint64, callId []types.U128) ([]*gtypes.TEECall, []types.U128, []types.StorageKey, error) {
// 	var pallet, method = "Bridge", "TEECalls"
// 	calls := make([]interface{}, 0, len(callId))
// 	for _, id := range callId {
// 		calls = append(calls, id)
// 	}
// 	set, err := m.QueryDoubleMapKeys(pallet, method, cid, calls, nil)
// 	if err != nil {
// 		return nil, nil, nil, err
// 	}

// 	var list []*gtypes.TEECall = make([]*gtypes.TEECall, 0, len(set))
// 	var keys []types.StorageKey = make([]types.StorageKey, 0, len(set))
// 	var callIds []types.U128 = make([]types.U128, 0, len(set))
// 	for _, elem := range set {
// 		for _, change := range elem.Changes {
// 			var d gtypes.TEECall

// 			if err := codec.Decode(change.StorageData, &d); err != nil && change.StorageData != nil {
// 				continue
// 			}
// 			keys = append(keys, change.StorageKey)
// 			list = append(list, &d)
// 			callIds = append(callIds, d.Id)
// 		}
// 	}

// 	return list, callIds, keys, nil
// }

// func (m *Chain) GetMetaApi(w gtypes.WorkId) (gtypes.ApiMeta, error) {
// 	api, ok, err := bridge.GetApiMetasLatest(m.Api.RPC.State, w)
// 	if err != nil {
// 		return gtypes.ApiMeta{}, err
// 	}
// 	if !ok {
// 		return gtypes.ApiMeta{}, fmt.Errorf("not found")
// 	}
// 	return api, nil
// }

// func (m *Chain) TeeCallback(cid uint64, callId []types.U128, callbacks []util.TeeCallBack) error {
// 	calls := make([]gtypes.RuntimeCall, 0, len(callbacks))
// 	for i, cb := range callbacks {
// 		var err []byte
// 		var isErr bool
// 		if cb.Err != "" {
// 			err = []byte(cb.Err)
// 			isErr = true
// 		}

// 		call := bridge.MakeInkCallbackCall(cid, callId[i], cb.Args, types.NewU128(*big.NewInt(0)), gtypes.OptionTByteSlice{
// 			IsSome:       isErr,
// 			IsNone:       !isErr,
// 			AsSomeField0: err,
// 		})
// 		calls = append(calls, call)
// 	}

// 	call := utility.MakeBatchCall(calls)
// 	return m.SignAndSubmit(m.signer, call, true)
// }
