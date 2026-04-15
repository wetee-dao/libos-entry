package cloud

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	chain "github.com/wetee-dao/ink.go"
	"github.com/wetee-dao/ink.go/util"
)

func DeployCloudWithNew(__ink_params chain.DeployParams) (*types.H160, error) {
	return __ink_params.Client.DeployContract(
		__ink_params.Code, __ink_params.Signer, types.NewU128(*big.NewInt(0)),
		util.InkContractInput{
			Selector: "0x00000000",
			Args:     []any{},
		},
		__ink_params.Salt,
	)
}

func InitCloudContract(client *chain.ChainClient, address string) (*Cloud, error) {
	contractAddress, err := util.HexToH160(address)
	if err != nil {
		return nil, err
	}
	return &Cloud{
		ChainClient: client,
		Address:     contractAddress,
	}, nil
}

type Cloud struct {
	ChainClient *chain.ChainClient
	Address     types.H160
}

func (c *Cloud) Client() *chain.ChainClient {
	return c.ChainClient
}

func (c *Cloud) ContractAddress() types.H160 {
	return c.Address
}

func (c *Cloud) DryRunInit(
	subnet_addr types.H160, pod_code_hash types.H256, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "init")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x25b9ac95",
			Args:     []any{subnet_addr, pod_code_hash},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecInit(
	subnet_addr types.H160, pod_code_hash types.H256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunInit(subnet_addr, pod_code_hash, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x25b9ac95",
			Args:     []any{subnet_addr, pod_code_hash},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfInit(
	subnet_addr types.H160, pod_code_hash types.H256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunInit(subnet_addr, pod_code_hash, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x25b9ac95",
			Args:     []any{subnet_addr, pod_code_hash},
		},
	)
}

func (c *Cloud) DryRunSetPodContract(
	pod_contract types.H256, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "set_pod_contract")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x0c28546b",
			Args:     []any{pod_contract},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecSetPodContract(
	pod_contract types.H256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetPodContract(pod_contract, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x0c28546b",
			Args:     []any{pod_contract},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfSetPodContract(
	pod_contract types.H256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSetPodContract(pod_contract, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x0c28546b",
			Args:     []any{pod_contract},
		},
	)
}

func (c *Cloud) QueryPodContract(
	__ink_params chain.DryRunParams,
) (*types.H256, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pod_contract")
	}
	v, gas, err := chain.DryRunInk[types.H256](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1417e3e6",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunUpdatePodContract(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "update_pod_contract")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x7b105159",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecUpdatePodContract(
	pod_id uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunUpdatePodContract(pod_id, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x7b105159",
			Args:     []any{pod_id},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfUpdatePodContract(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunUpdatePodContract(pod_id, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x7b105159",
			Args:     []any{pod_id},
		},
	)
}

func (c *Cloud) DryRunSetMintInterval(
	t uint32, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "set_mint_interval")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb3b8cfcf",
			Args:     []any{t},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecSetMintInterval(
	t uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunSetMintInterval(t, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xb3b8cfcf",
			Args:     []any{t},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfSetMintInterval(
	t uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunSetMintInterval(t, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xb3b8cfcf",
			Args:     []any{t},
		},
	)
}

func (c *Cloud) QueryMintInterval(
	__ink_params chain.DryRunParams,
) (*uint32, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "mint_interval")
	}
	v, gas, err := chain.DryRunInk[uint32](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x6f7ded81",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QuerySubnetAddress(
	__ink_params chain.DryRunParams,
) (*types.H160, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "subnet_address")
	}
	v, gas, err := chain.DryRunInk[types.H160](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x9ddb90ba",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryCharge(
	__ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "charge")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x96dc00c1",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) QueryBalance(
	asset AssetInfo, __ink_params chain.DryRunParams,
) (*types.U256, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "balance")
	}
	v, gas, err := chain.DryRunInk[types.U256](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x017f8fa1",
			Args:     []any{asset},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryPodLen(
	__ink_params chain.DryRunParams,
) (*uint64, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pod_len")
	}
	v, gas, err := chain.DryRunInk[uint64](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x9b8bc26c",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryPods(
	start util.Option[uint64], size uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_36, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pods")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_36](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x77ffaeb8",
			Args:     []any{start, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryUserPodLen(
	__ink_params chain.DryRunParams,
) (*uint64, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "user_pod_len")
	}
	v, gas, err := chain.DryRunInk[uint64](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb2624c68",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryUserPods(
	start util.Option[uint64], size uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_36, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "user_pods")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_36](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x1474c16c",
			Args:     []any{start, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryWorkerPodsVersion(
	worker_id uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_39, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "worker_pods_version")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_39](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xe2b186e9",
			Args:     []any{worker_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryWorkerPods(
	worker_id uint64, start util.Option[uint64], size uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_36, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "worker_pods")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_36](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd13982c0",
			Args:     []any{worker_id, start, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryWorkerPodLen(
	worker_id uint64, __ink_params chain.DryRunParams,
) (*uint64, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "worker_pod_len")
	}
	v, gas, err := chain.DryRunInk[uint64](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xa64e8709",
			Args:     []any{worker_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryUserSecrets(
	user types.H160, start util.Option[uint64], size uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_44, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "user_secrets")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_44](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x09b8ee81",
			Args:     []any{user, start, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QuerySecret(
	user types.H160, index uint64, __ink_params chain.DryRunParams,
) (*util.Option[Secret], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "secret")
	}
	v, gas, err := chain.DryRunInk[util.Option[Secret]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x66e75470",
			Args:     []any{user, index},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunCreateSecret(
	key []byte, hash types.H256, __ink_params chain.DryRunParams,
) (*util.Result[uint64, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "create_secret")
	}
	v, gas, err := chain.DryRunInk[util.Result[uint64, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xc8f2dd2a",
			Args:     []any{key, hash},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecCreateSecret(
	key []byte, hash types.H256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCreateSecret(key, hash, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xc8f2dd2a",
			Args:     []any{key, hash},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfCreateSecret(
	key []byte, hash types.H256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunCreateSecret(key, hash, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xc8f2dd2a",
			Args:     []any{key, hash},
		},
	)
}

func (c *Cloud) DryRunMintSecret(
	user types.H160, index uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "mint_secret")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x77993afd",
			Args:     []any{user, index},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecMintSecret(
	user types.H160, index uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunMintSecret(user, index, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x77993afd",
			Args:     []any{user, index},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfMintSecret(
	user types.H160, index uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunMintSecret(user, index, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x77993afd",
			Args:     []any{user, index},
		},
	)
}

func (c *Cloud) DryRunDelSecret(
	index uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "del_secret")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xe71b690d",
			Args:     []any{index},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecDelSecret(
	index uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunDelSecret(index, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xe71b690d",
			Args:     []any{index},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfDelSecret(
	index uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunDelSecret(index, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xe71b690d",
			Args:     []any{index},
		},
	)
}

func (c *Cloud) DryRunCreateDisk(
	key []byte, size uint32, __ink_params chain.DryRunParams,
) (*util.Result[uint64, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "create_disk")
	}
	v, gas, err := chain.DryRunInk[util.Result[uint64, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xe1db4ae6",
			Args:     []any{key, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecCreateDisk(
	key []byte, size uint32, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCreateDisk(key, size, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xe1db4ae6",
			Args:     []any{key, size},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfCreateDisk(
	key []byte, size uint32, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunCreateDisk(key, size, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xe1db4ae6",
			Args:     []any{key, size},
		},
	)
}

func (c *Cloud) DryRunUpdateDiskKey(
	user types.H160, id uint64, hash types.H256, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "update_disk_key")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xb951dcf3",
			Args:     []any{user, id, hash},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecUpdateDiskKey(
	user types.H160, id uint64, hash types.H256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunUpdateDiskKey(user, id, hash, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xb951dcf3",
			Args:     []any{user, id, hash},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfUpdateDiskKey(
	user types.H160, id uint64, hash types.H256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunUpdateDiskKey(user, id, hash, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xb951dcf3",
			Args:     []any{user, id, hash},
		},
	)
}

func (c *Cloud) QueryDisk(
	user types.H160, disk_id uint64, __ink_params chain.DryRunParams,
) (*util.Option[Disk], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "disk")
	}
	v, gas, err := chain.DryRunInk[util.Option[Disk]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x69e62161",
			Args:     []any{user, disk_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryUserDisks(
	user types.H160, start util.Option[uint64], size uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_55, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "user_disks")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_55](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x77c9ef5d",
			Args:     []any{user, start, size},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunDelDisk(
	disk_id uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "del_disk")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x5d82148f",
			Args:     []any{disk_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecDelDisk(
	disk_id uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunDelDisk(disk_id, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x5d82148f",
			Args:     []any{disk_id},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfDelDisk(
	disk_id uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunDelDisk(disk_id, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x5d82148f",
			Args:     []any{disk_id},
		},
	)
}

func (c *Cloud) QueryPodExtInfo(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*util.Option[Tuple_66], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pod_ext_info")
	}
	v, gas, err := chain.DryRunInk[util.Option[Tuple_66]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x0fdf6b91",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) QueryPodsByIds(
	pod_ids []uint64, __ink_params chain.DryRunParams,
) (*[]Tuple_74, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pods_by_ids")
	}
	v, gas, err := chain.DryRunInk[[]Tuple_74](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x8a276837",
			Args:     []any{pod_ids},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunTransfer(
	asset AssetInfo, to types.H160, amount types.U256, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "transfer")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd2895918",
			Args:     []any{asset, to, amount},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecTransfer(
	asset AssetInfo, to types.H160, amount types.U256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunTransfer(asset, to, amount, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xd2895918",
			Args:     []any{asset, to, amount},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfTransfer(
	asset AssetInfo, to types.H160, amount types.U256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunTransfer(asset, to, amount, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0xd2895918",
			Args:     []any{asset, to, amount},
		},
	)
}

func (c *Cloud) QueryPod(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*util.Option[Tuple_77], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pod")
	}
	v, gas, err := chain.DryRunInk[util.Option[Tuple_77]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xefc61efb",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunCreatePod(
	name []byte, pod_type PodType, tee_type TEEType, containers []Container, region_id uint32, level byte, pay_asset uint32, worker_id uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "create_pod")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x7d7d3b36",
			Args:     []any{name, pod_type, tee_type, containers, region_id, level, pay_asset, worker_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecCreatePod(
	name []byte, pod_type PodType, tee_type TEEType, containers []Container, region_id uint32, level byte, pay_asset uint32, worker_id uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunCreatePod(name, pod_type, tee_type, containers, region_id, level, pay_asset, worker_id, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x7d7d3b36",
			Args:     []any{name, pod_type, tee_type, containers, region_id, level, pay_asset, worker_id},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfCreatePod(
	name []byte, pod_type PodType, tee_type TEEType, containers []Container, region_id uint32, level byte, pay_asset uint32, worker_id uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunCreatePod(name, pod_type, tee_type, containers, region_id, level, pay_asset, worker_id, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x7d7d3b36",
			Args:     []any{name, pod_type, tee_type, containers, region_id, level, pay_asset, worker_id},
		},
	)
}

func (c *Cloud) DryRunStartPod(
	pod_id uint64, pod_key util.AccountId, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "start_pod")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x54ac72d2",
			Args:     []any{pod_id, pod_key},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecStartPod(
	pod_id uint64, pod_key util.AccountId, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunStartPod(pod_id, pod_key, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x54ac72d2",
			Args:     []any{pod_id, pod_key},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfStartPod(
	pod_id uint64, pod_key util.AccountId, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunStartPod(pod_id, pod_key, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x54ac72d2",
			Args:     []any{pod_id, pod_key},
		},
	)
}

func (c *Cloud) DryRunStopPod(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "stop_pod")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x3c644b83",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecStopPod(
	pod_id uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunStopPod(pod_id, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x3c644b83",
			Args:     []any{pod_id},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfStopPod(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunStopPod(pod_id, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x3c644b83",
			Args:     []any{pod_id},
		},
	)
}

func (c *Cloud) DryRunRestartPod(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "restart_pod")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x9b6b5d51",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecRestartPod(
	pod_id uint64, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunRestartPod(pod_id, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x9b6b5d51",
			Args:     []any{pod_id},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfRestartPod(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunRestartPod(pod_id, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x9b6b5d51",
			Args:     []any{pod_id},
		},
	)
}

func (c *Cloud) DryRunMintPod(
	pod_id uint64, report types.H256, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "mint_pod")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x9068edc7",
			Args:     []any{pod_id, report},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecMintPod(
	pod_id uint64, report types.H256, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunMintPod(pod_id, report, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x9068edc7",
			Args:     []any{pod_id, report},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfMintPod(
	pod_id uint64, report types.H256, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunMintPod(pod_id, report, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x9068edc7",
			Args:     []any{pod_id, report},
		},
	)
}

func (c *Cloud) QueryPodReport(
	pod_id uint64, __ink_params chain.DryRunParams,
) (*util.Option[types.H256], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "pod_report")
	}
	v, gas, err := chain.DryRunInk[util.Option[types.H256]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd6901a3f",
			Args:     []any{pod_id},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}

func (c *Cloud) DryRunEditContainer(
	pod_id uint64, containers []ContainerInput, __ink_params chain.DryRunParams,
) (*util.Result[util.NullTuple, Error], *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "edit_container")
	}
	v, gas, err := chain.DryRunInk[util.Result[util.NullTuple, Error]](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0x52051579",
			Args:     []any{pod_id, containers},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	if v != nil && v.IsErr {
		return nil, nil, errors.New("Contract Reverted: " + v.E.Error())
	}

	return v, gas, nil
}

func (c *Cloud) ExecEditContainer(
	pod_id uint64, containers []ContainerInput, __ink_params chain.ExecParams,
) error {
	_param := chain.DefaultParamWithOrigin(__ink_params.Signer.AccountID())
	_param.PayAmount = __ink_params.PayAmount
	_, gas, err := c.DryRunEditContainer(pod_id, containers, _param)
	if err != nil {
		return err
	}
	return chain.CallInk(
		c,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x52051579",
			Args:     []any{pod_id, containers},
		},
		__ink_params,
	)
}

func (c *Cloud) CallOfEditContainer(
	pod_id uint64, containers []ContainerInput, __ink_params chain.DryRunParams,
) (*types.Call, error) {
	_, gas, err := c.DryRunEditContainer(pod_id, containers, __ink_params)
	if err != nil {
		return nil, err
	}
	return chain.CallOfTransaction(
		c,
		__ink_params.PayAmount,
		gas.GasRequired,
		gas.StorageDeposit,
		util.InkContractInput{
			Selector: "0x52051579",
			Args:     []any{pod_id, containers},
		},
	)
}

func (c *Cloud) QuerySubnetSideChainKey(
	__ink_params chain.DryRunParams,
) (*types.H160, *chain.DryRunReturnGas, error) {
	if c.ChainClient.Debug {
		fmt.Println()
		util.LogWithPurple("[ DryRun   method ]", "subnet_side_chain_key")
	}
	v, gas, err := chain.DryRunInk[types.H160](
		c,
		__ink_params.Origin,
		__ink_params.PayAmount,
		__ink_params.GasLimit,
		__ink_params.StorageDepositLimit,
		util.InkContractInput{
			Selector: "0xd99d60f8",
			Args:     []any{},
		},
	)
	if err != nil && !errors.Is(err, chain.ErrContractReverted) {
		return nil, nil, err
	}
	return v, gas, nil
}
