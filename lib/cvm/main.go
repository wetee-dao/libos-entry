package ego

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/wetee-dao/go-sdk/core"
	"github.com/wetee-dao/libos-entry/libos"
	"github.com/wetee-dao/libos-entry/util"
)

func InitEgo() error {
	hostfs := &CvmFs{}
	return libos.PreLoad(hostfs)
}

type CvmFs struct {
	afero.OsFs
	report     []byte
	lastReport int64
}

// Read implements util.Fs.
func (e *CvmFs) ReadFile(filename string) ([]byte, error) {
	bt, err := afero.ReadFile(e, filename)
	if err != nil {
		return nil, err
	}

	// 解密数据
	keyBytes, err := e.Decrypt(bt)
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

// Write implements util.Fs.
func (e *CvmFs) WriteFile(filename string, data []byte, perm os.FileMode) error {

	// 加密数据
	val, err := e.Encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt Key: %v", err)
	}

	return afero.WriteFile(e, filename, val, perm)
}

// Decrypt implements libos.SecretFunction.
func (e *CvmFs) Decrypt(val []byte) ([]byte, error) {
	return val, nil
}

// Encrypt implements libos.SecretFunction.
func (e *CvmFs) Encrypt(val []byte) ([]byte, error) {
	return val, nil
}

func (i *CvmFs) IssueReport(pk *core.Signer, data []byte) (*util.TeeParam, error) {
	timestamp := time.Now().Unix()
	if i.report != nil && i.lastReport+30 > timestamp {
		return &util.TeeParam{
			Time:    i.lastReport,
			Address: pk.SS58Address(42),
			Report:  i.report,
			Data:    data,
		}, nil
	}

	return &util.TeeParam{
		Time:    timestamp,
		Address: pk.SS58Address(42),
		Report:  []byte{},
		Data:    data,
	}, nil
}

func (e *CvmFs) VerifyReport(workerReport *util.TeeParam) (*util.TeeReport, error) {
	return &util.TeeReport{
		CodeSignature: []byte{},
		CodeSigner:    []byte{},
		CodeProductID: []byte{},
	}, nil
}
