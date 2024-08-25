package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/go-sdk/pallet/types"
)

// 定义返回的数据结构
type TeeCallBack struct {
	Err  string         `json:"err"`
	Args []types.InkArg `json:"args"`
}

type TeeCallBackJSON struct {
	Err  string   `json:"err"`
	Args []string `json:"args"`
}

func (t *TeeCallBack) ToJSON() ([]byte, error) {
	args := make([]string, 0, len(t.Args))
	for _, arg := range t.Args {
		bt, _ := codec.Encode(arg)
		args = append(args, base64.StdEncoding.EncodeToString(bt))
	}

	data := TeeCallBackJSON{
		Err:  t.Err,
		Args: args,
	}

	return json.Marshal(data)
}

func UnmarshalToArgs(bt []byte) *TeeCallBack {
	// 解析 JSON 数据
	var data TeeCallBackJSON
	err := json.Unmarshal(bt, &data)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return nil
	}

	// 遍历解析后的数组
	args := make([]types.InkArg, 0, len(data.Args))
	for _, str := range data.Args {
		var arg types.InkArg
		argStr, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil
		}
		err = codec.Decode(argStr, &arg)
		if err != nil {
			return nil
		}
		args = append(args, arg)
	}

	return &TeeCallBack{
		Err:  data.Err,
		Args: args,
	}
}

// func transferArgs(args []map[string]interface{}) []types.InkArg {
// 	var args2 = make([]types.InkArg, 0, len(args))
// 	for i, arg := range args {
// 		for k, v := range arg {
// 			switch v {
// 			case "InkArg::Bool":
// 				args2 = append(args2, types.InkArg{
// 					IsBool:       true,
// 					AsBoolField0: v.(bool),
// 				})
// 				break
// 			case "InkArg::U8":
// 				args2 = append(args2, types.InkArg{
// 					IsU8:       true,
// 					AsU8Field0: v.(uint8),
// 				})
// 				break
// 			case "InkArg::I8":
// 				args2 = append(args2, types.InkArg{
// 					IsI8:       true,
// 					AsI8Field0: v.(int8),
// 				})
// 				break
// 			case "InkArg::U16":
// 				args2 = append(args2, types.InkArg{
// 					IsU16:       true,
// 					AsU16Field0: v.(uint16),
// 				})
// 				break
// 			case "InkArg::I16":
// 				args2 = append(args2, types.InkArg{
// 					IsI16:       true,
// 					AsI16Field0: v.(int16),
// 				})
// 				break
// 			case "InkArg::U32":
// 				args2 = append(args2, types.InkArg{
// 					IsU32:       true,
// 					AsU32Field0: v.(uint32),
// 				})
// 				break
// 			case "InkArg::I32":
// 				args2 = append(args2, types.InkArg{
// 					IsI32:       true,
// 					AsI32Field0: v.(int32),
// 				})
// 				break
// 			case "InkArg::U64":
// 				args2 = append(args2, types.InkArg{
// 					IsU64:       true,
// 					AsU64Field0: v.(uint64),
// 				})
// 				break
// 			case "InkArg::I64":
// 				args2 = append(args2, types.InkArg{
// 					IsI64:       true,
// 					AsI64Field0: v.(int64),
// 				})
// 				break
// 			case "InkArg::U128":
// 				i := stypes.NewU128(*big.NewInt(v.(int64)))
// 				args2 = append(args2, types.InkArg{
// 					IsU128:       true,
// 					AsU128Field0: i,
// 				})
// 				break
// 			case "InkArg::I128":
// 				args2 = append(args2, types.InkArg{
// 					IsI128:       true,
// 					AsI128Field0: v.(int64),
// 				})
// 				break
// 			case "InkArg::TString":
// 				bt, _ := base64.StdEncoding.DecodeString(v.(string))
// 				args2 = append(args2, types.InkArg{
// 					IsTString:       true,
// 					AsTStringField0: bt,
// 				})
// 				break
// 			default:
// 				fmt.Println("transferArgs", "err", "unknown type", k, v)
// 				break
// 			}
// 		}
// 	}

// }

// func (ty InkArg) MarshalJSON() ([]byte, error) {
// 	if ty.IsBool {
// 		m := map[string]interface{}{"InkArg::Bool": ty.AsBoolField0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsU8 {
// 		m := map[string]interface{}{"InkArg::U8": ty.AsU8Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsI8 {
// 		m := map[string]interface{}{"InkArg::I8": ty.AsI8Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsU16 {
// 		m := map[string]interface{}{"InkArg::U16": ty.AsU16Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsI16 {
// 		m := map[string]interface{}{"InkArg::I16": ty.AsI16Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsU32 {
// 		m := map[string]interface{}{"InkArg::U32": ty.AsU32Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsI32 {
// 		m := map[string]interface{}{"InkArg::I32": ty.AsI32Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsU64 {
// 		m := map[string]interface{}{"InkArg::U64": ty.AsU64Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsI64 {
// 		m := map[string]interface{}{"InkArg::I64": ty.AsI64Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsU128 {
// 		m := map[string]interface{}{"InkArg::U128": ty.AsU128Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsI128 {
// 		m := map[string]interface{}{"InkArg::I128": ty.AsI128Field0}
// 		return json.Marshal(m)
// 	}
// 	if ty.IsTString {
// 		m := map[string]interface{}{"InkArg::TString": ty.AsTStringField0}
// 		return json.Marshal(m)
// 	}
// 	return nil, fmt.Errorf("No variant detected")
// }
