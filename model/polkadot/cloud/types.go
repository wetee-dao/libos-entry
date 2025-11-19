package cloud

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/ink.go/util"
)

type Pod struct { // Composite
	Name       []byte
	Owner      types.H160
	Contract   types.H160
	Ptype      PodType
	StartBlock uint32
	TeeType    TEEType
}
type PodType struct { // Enum
	CPU    *bool // 0
	GPU    *bool // 1
	SCRIPT *bool // 2
}

func (ty PodType) Encode(encoder scale.Encoder) (err error) {
	if ty.CPU != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.GPU != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.SCRIPT != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *PodType) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.CPU = &t
		return
	case 1: // Base
		t := true
		ty.GPU = &t
		return
	case 2: // Base
		t := true
		ty.SCRIPT = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type TEEType struct { // Enum
	SGX *bool // 0
	CVM *bool // 1
}

func (ty TEEType) Encode(encoder scale.Encoder) (err error) {
	if ty.SGX != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.CVM != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *TEEType) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.SGX = &t
		return
	case 1: // Base
		t := true
		ty.CVM = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Service struct { // Enum
	Tcp        *uint16 // 0
	Udp        *uint16 // 1
	Http       *uint16 // 2
	Https      *uint16 // 3
	ProjectTcp *uint16 // 4
	ProjectUdp *uint16 // 5
}

func (ty Service) Encode(encoder scale.Encoder) (err error) {
	if ty.Tcp != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Tcp)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Udp != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Udp)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Http != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Http)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.Https != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.Https)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ProjectTcp != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.ProjectTcp)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ProjectUdp != nil {
		err = encoder.PushByte(5)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.ProjectUdp)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Service) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Inline
		ty.Tcp = new(uint16)
		err = decoder.Decode(ty.Tcp)
		if err != nil {
			return err
		}
		return
	case 1: // Inline
		ty.Udp = new(uint16)
		err = decoder.Decode(ty.Udp)
		if err != nil {
			return err
		}
		return
	case 2: // Inline
		ty.Http = new(uint16)
		err = decoder.Decode(ty.Http)
		if err != nil {
			return err
		}
		return
	case 3: // Inline
		ty.Https = new(uint16)
		err = decoder.Decode(ty.Https)
		if err != nil {
			return err
		}
		return
	case 4: // Inline
		ty.ProjectTcp = new(uint16)
		err = decoder.Decode(ty.ProjectTcp)
		if err != nil {
			return err
		}
		return
	case 5: // Inline
		ty.ProjectUdp = new(uint16)
		err = decoder.Decode(ty.ProjectUdp)
		if err != nil {
			return err
		}
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type ContainerDisk struct { // Composite
	Id   uint64
	Path []byte
}
type Env struct { // Enum
	Env *struct { // 0
		F0 []byte
		F1 []byte
	}
	File *struct { // 1
		F0 []byte
		F1 []byte
	}
	Encrypt *struct { // 2
		F0 []byte
		F1 uint64
	}
}

func (ty Env) Encode(encoder scale.Encoder) (err error) {
	if ty.Env != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.Env.F0)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.Env.F1)
		if err != nil {
			return err
		}

		return nil
	}

	if ty.File != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.File.F0)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.File.F1)
		if err != nil {
			return err
		}

		return nil
	}

	if ty.Encrypt != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.Encrypt.F0)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.Encrypt.F1)
		if err != nil {
			return err
		}

		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Env) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Tuple
		ty.Env = &struct {
			F0 []byte
			F1 []byte
		}{}

		err = decoder.Decode(&ty.Env.F0)
		if err != nil {
			return err
		}

		err = decoder.Decode(&ty.Env.F1)
		if err != nil {
			return err
		}

		return
	case 1: // Tuple
		ty.File = &struct {
			F0 []byte
			F1 []byte
		}{}

		err = decoder.Decode(&ty.File.F0)
		if err != nil {
			return err
		}

		err = decoder.Decode(&ty.File.F1)
		if err != nil {
			return err
		}

		return
	case 2: // Tuple
		ty.Encrypt = &struct {
			F0 []byte
			F1 uint64
		}{}

		err = decoder.Decode(&ty.Encrypt.F0)
		if err != nil {
			return err
		}

		err = decoder.Decode(&ty.Encrypt.F1)
		if err != nil {
			return err
		}

		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Container struct { // Composite
	Image   []byte
	Command Command
	Port    []Service
	Cpu     uint32
	Mem     uint32
	Disk    []ContainerDisk
	Gpu     uint32
	Env     []Env
}
type Command struct { // Enum
	SH   *[]byte // 0
	BASH *[]byte // 1
	ZSH  *[]byte // 2
	NONE *bool   // 3
}

func (ty Command) Encode(encoder scale.Encoder) (err error) {
	if ty.SH != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.SH)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.BASH != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.BASH)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.ZSH != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.ZSH)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NONE != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Command) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Inline
		ty.SH = new([]byte)
		err = decoder.Decode(ty.SH)
		if err != nil {
			return err
		}
		return
	case 1: // Inline
		ty.BASH = new([]byte)
		err = decoder.Decode(ty.BASH)
		if err != nil {
			return err
		}
		return
	case 2: // Inline
		ty.ZSH = new([]byte)
		err = decoder.Decode(ty.ZSH)
		if err != nil {
			return err
		}
		return
	case 3: // Base
		t := true
		ty.NONE = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Secret struct { // Composite
	K      []byte
	Hash   types.H256
	Minted bool
}
type Disk struct { // Enum
	SecretSSD *struct { // 0
		F0 []byte
		F1 []byte
		F2 uint32
	}
}

func (ty Disk) Encode(encoder scale.Encoder) (err error) {
	if ty.SecretSSD != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.SecretSSD.F0)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.SecretSSD.F1)
		if err != nil {
			return err
		}

		err = encoder.Encode(ty.SecretSSD.F2)
		if err != nil {
			return err
		}

		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Disk) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Tuple
		ty.SecretSSD = &struct {
			F0 []byte
			F1 []byte
			F2 uint32
		}{}

		err = decoder.Decode(&ty.SecretSSD.F0)
		if err != nil {
			return err
		}

		err = decoder.Decode(&ty.SecretSSD.F1)
		if err != nil {
			return err
		}

		err = decoder.Decode(&ty.SecretSSD.F2)
		if err != nil {
			return err
		}

		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Error struct { // Enum
	SetCodeFailed          *bool // 0
	MustCallByGovContract  *bool // 1
	WorkerLevelNotEnough   *bool // 2
	RegionNotMatch         *bool // 3
	WorkerNotOnline        *bool // 4
	NotPodOwner            *bool // 5
	PodKeyNotExist         *bool // 6
	PodStatusError         *bool // 7
	InvalidSideChainCaller *bool // 8
	DelFailed              *bool // 9
	NotFound               *bool // 10
}

func (ty Error) Encode(encoder scale.Encoder) (err error) {
	if ty.SetCodeFailed != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.MustCallByGovContract != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerLevelNotEnough != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.RegionNotMatch != nil {
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.WorkerNotOnline != nil {
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NotPodOwner != nil {
		err = encoder.PushByte(5)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.PodKeyNotExist != nil {
		err = encoder.PushByte(6)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.PodStatusError != nil {
		err = encoder.PushByte(7)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.InvalidSideChainCaller != nil {
		err = encoder.PushByte(8)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.DelFailed != nil {
		err = encoder.PushByte(9)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.NotFound != nil {
		err = encoder.PushByte(10)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *Error) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.SetCodeFailed = &t
		return
	case 1: // Base
		t := true
		ty.MustCallByGovContract = &t
		return
	case 2: // Base
		t := true
		ty.WorkerLevelNotEnough = &t
		return
	case 3: // Base
		t := true
		ty.RegionNotMatch = &t
		return
	case 4: // Base
		t := true
		ty.WorkerNotOnline = &t
		return
	case 5: // Base
		t := true
		ty.NotPodOwner = &t
		return
	case 6: // Base
		t := true
		ty.PodKeyNotExist = &t
		return
	case 7: // Base
		t := true
		ty.PodStatusError = &t
		return
	case 8: // Base
		t := true
		ty.InvalidSideChainCaller = &t
		return
	case 9: // Base
		t := true
		ty.DelFailed = &t
		return
	case 10: // Base
		t := true
		ty.NotFound = &t
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}
func (ty *Error) Error() string {
	if ty.SetCodeFailed != nil {
		return "SetCodeFailed"
	}

	if ty.MustCallByGovContract != nil {
		return "MustCallByGovContract"
	}

	if ty.WorkerLevelNotEnough != nil {
		return "WorkerLevelNotEnough"
	}

	if ty.RegionNotMatch != nil {
		return "RegionNotMatch"
	}

	if ty.WorkerNotOnline != nil {
		return "WorkerNotOnline"
	}

	if ty.NotPodOwner != nil {
		return "NotPodOwner"
	}

	if ty.PodKeyNotExist != nil {
		return "PodKeyNotExist"
	}

	if ty.PodStatusError != nil {
		return "PodStatusError"
	}

	if ty.InvalidSideChainCaller != nil {
		return "InvalidSideChainCaller"
	}

	if ty.DelFailed != nil {
		return "DelFailed"
	}

	if ty.NotFound != nil {
		return "NotFound"
	}
	return "Unknown"
}

type ContainerInput struct { // Composite
	Etype     EditType
	Container Container
}
type EditType struct { // Enum
	INSERT *bool   // 0
	UPDATE *uint64 // 1
	REMOVE *uint64 // 2
}

func (ty EditType) Encode(encoder scale.Encoder) (err error) {
	if ty.INSERT != nil {
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.UPDATE != nil {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.UPDATE)
		if err != nil {
			return err
		}
		return nil
	}

	if ty.REMOVE != nil {
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(*ty.REMOVE)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unrecognized enum")
}

func (ty *EditType) Decode(decoder scale.Decoder) (err error) {
	variant, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch variant {
	case 0: // Base
		t := true
		ty.INSERT = &t
		return
	case 1: // Inline
		ty.UPDATE = new(uint64)
		err = decoder.Decode(ty.UPDATE)
		if err != nil {
			return err
		}
		return
	case 2: // Inline
		ty.REMOVE = new(uint64)
		err = decoder.Decode(ty.REMOVE)
		if err != nil {
			return err
		}
		return
	default:
		return fmt.Errorf("unrecognized enum")
	}
}

type Tuple_115 struct { // Tuple
	F0 uint64
	F1 Pod
	F2 []Tuple_117
	F3 byte
}
type Tuple_117 struct { // Tuple
	F0 uint64
	F1 Container
}
type Tuple_121 struct { // Tuple
	F0 uint64
	F1 uint32
	F2 uint32
	F3 byte
}
type Tuple_124 struct { // Tuple
	F0 Pod
	F1 []Tuple_117
	F2 uint32
	F3 byte
}
type Tuple_127 struct { // Tuple
	F0 uint64
	F1 K8sCluster
	F2 []byte
}
type K8sCluster struct { // Composite
	Name          []byte
	Owner         types.H160
	Level         byte
	RegionId      uint32
	StartBlock    uint32
	StopBlock     util.Option[uint32]
	TerminalBlock util.Option[uint32]
	P2pId         util.AccountId
	Ip            Ip
	Port          uint32
	Status        byte
}
type Ip struct { // Composite
	Ipv4   util.Option[uint32]
	Ipv6   util.Option[types.U128]
	Domain util.Option[[]byte]
}
type Tuple_135 struct { // Tuple
	F0 uint64
	F1 Pod
	F2 []Tuple_137
	F3 uint32
	F4 uint32
	F5 byte
}
type Tuple_137 struct { // Tuple
	F0 uint64
	F1 Tuple_138
}
type Tuple_138 struct { // Tuple
	F0 Container
	F1 []util.Option[Disk]
}
type Tuple_143 struct { // Tuple
	F0 uint64
	F1 Secret
}
type Tuple_151 struct { // Tuple
	F0 uint64
	F1 Disk
}
