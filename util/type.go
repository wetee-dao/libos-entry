package util

type LoadParam struct {
	Address   string
	Time      string
	Signature string
	Cert      []byte
	Report    []byte
}

type Secrets struct {
	Files map[string]string
	Env   map[string]string
}

type SecretFunction interface {
	VerifyReport(reportBytes, certBytes, signer []byte) error
	Encrypt(val []byte) ([]byte, error)
	Decrypt(val []byte) ([]byte, error)
}
