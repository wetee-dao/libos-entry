package model

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
// 	"github.com/gogo/protobuf/proto"
// )

// // 与 encodePayload / decodePayload 中 switch 顺序一致。
// const (
// 	txVarEmpty uint8 = iota
// 	txVarEpochEnd
// 	txVarEpochStart
// 	txVarHubCall
// 	txVarSyncTxStart
// 	txVarSyncTxEnd
// 	txVarSyncTxRetry
// 	txVarContract
// )

// // txWire 为 Tx 的 SCALE 线格式：载荷分支 + 内层字节 + 元数据；验签时 Signature 为空切片。
// type txWire struct {
// 	Variant       uint8
// 	Inner         []byte
// 	Caller        []byte
// 	Signature     []byte
// 	SignatureType uint32
// }

// func encodePayload(p isTx_Payload) (variant uint8, inner []byte, err error) {
// 	switch t := p.(type) {
// 	case *Tx_Empty:
// 		b, err := codec.Encode(t.Empty)
// 		return txVarEmpty, b, err
// 	case *Tx_EpochEnd:
// 		if t.EpochEnd == nil {
// 			return 0, nil, errors.New("tx: epoch_end is nil")
// 		}
// 		b, err := proto.Marshal(t.EpochEnd)
// 		return txVarEpochEnd, b, err
// 	case *Tx_EpochStart:
// 		b, err := codec.Encode(t.EpochStart)
// 		return txVarEpochStart, b, err
// 	case *Tx_HubCall:
// 		if t.HubCall == nil {
// 			return 0, nil, errors.New("tx: hub_call is nil")
// 		}
// 		b, err := proto.Marshal(t.HubCall)
// 		return txVarHubCall, b, err
// 	case *Tx_SyncTxStart:
// 		b, err := codec.Encode(t.SyncTxStart)
// 		return txVarSyncTxStart, b, err
// 	case *Tx_SyncTxEnd:
// 		b, err := codec.Encode(t.SyncTxEnd)
// 		return txVarSyncTxEnd, b, err
// 	case *Tx_SyncTxRetry:
// 		b, err := codec.Encode(t.SyncTxRetry)
// 		return txVarSyncTxRetry, b, err
// 	case *Tx_Contract:
// 		if len(t.Contract) == 0 {
// 			return 0, nil, errors.New("tx: empty contract payload")
// 		}
// 		return txVarContract, append([]byte(nil), t.Contract...), nil
// 	default:
// 		return 0, nil, fmt.Errorf("tx: unknown payload type %T", p)
// 	}
// }

// func decodePayload(variant uint8, inner []byte) (isTx_Payload, error) {
// 	switch variant {
// 	case txVarEmpty:
// 		var v int64
// 		if err := codec.Decode(inner, &v); err != nil {
// 			return nil, fmt.Errorf("tx: decode empty: %w", err)
// 		}
// 		return &Tx_Empty{Empty: v}, nil
// 	case txVarEpochEnd:
// 		var e EpochEnd
// 		if err := proto.Unmarshal(inner, &e); err != nil {
// 			return nil, fmt.Errorf("tx: decode epoch_end: %w", err)
// 		}
// 		return &Tx_EpochEnd{EpochEnd: &e}, nil
// 	case txVarEpochStart:
// 		var v int64
// 		if err := codec.Decode(inner, &v); err != nil {
// 			return nil, fmt.Errorf("tx: decode epoch_start: %w", err)
// 		}
// 		return &Tx_EpochStart{EpochStart: v}, nil
// 	case txVarHubCall:
// 		var h HubCall
// 		if err := proto.Unmarshal(inner, &h); err != nil {
// 			return nil, fmt.Errorf("tx: decode hub_call: %w", err)
// 		}
// 		return &Tx_HubCall{HubCall: &h}, nil
// 	case txVarSyncTxStart:
// 		var v int64
// 		if err := codec.Decode(inner, &v); err != nil {
// 			return nil, fmt.Errorf("tx: decode sync_tx_start: %w", err)
// 		}
// 		return &Tx_SyncTxStart{SyncTxStart: v}, nil
// 	case txVarSyncTxEnd:
// 		var v int64
// 		if err := codec.Decode(inner, &v); err != nil {
// 			return nil, fmt.Errorf("tx: decode sync_tx_end: %w", err)
// 		}
// 		return &Tx_SyncTxEnd{SyncTxEnd: v}, nil
// 	case txVarSyncTxRetry:
// 		var v int64
// 		if err := codec.Decode(inner, &v); err != nil {
// 			return nil, fmt.Errorf("tx: decode sync_tx_retry: %w", err)
// 		}
// 		return &Tx_SyncTxRetry{SyncTxRetry: v}, nil
// 	case txVarContract:
// 		if len(inner) == 0 {
// 			return nil, errors.New("tx: empty contract payload")
// 		}
// 		return &Tx_Contract{Contract: append([]byte(nil), inner...)}, nil
// 	default:
// 		return nil, fmt.Errorf("tx: unknown variant %d", variant)
// 	}
// }

// // EncodeTx 将 Tx 编码为 SCALE 字节（用于 TxBox.tx 及网络传输）。
// func EncodeTx(tx *Tx) ([]byte, error) {
// 	if tx == nil {
// 		return nil, errors.New("tx is nil")
// 	}
// 	if tx.Payload == nil {
// 		return nil, errors.New("tx: missing payload")
// 	}
// 	v, inner, err := encodePayload(tx.Payload)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return codec.Encode(txWire{
// 		Variant:       v,
// 		Inner:         inner,
// 		Caller:        tx.Caller,
// 		Signature:     tx.Signature,
// 		SignatureType: tx.SignatureType,
// 	})
// }

// // DecodeTx 解析 EncodeTx 的输出。
// func DecodeTx(data []byte) (*Tx, error) {
// 	var w txWire
// 	if err := codec.Decode(data, &w); err != nil {
// 		return nil, err
// 	}
// 	p, err := decodePayload(w.Variant, w.Inner)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Tx{
// 		Payload:       p,
// 		Caller:        w.Caller,
// 		Signature:     w.Signature,
// 		SignatureType: w.SignatureType,
// 	}, nil
// }
