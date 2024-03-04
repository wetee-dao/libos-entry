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
