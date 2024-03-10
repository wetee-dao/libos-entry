package util

// 请求开始TEE容器
type LoadParam struct {
	Address   string
	Time      string
	Signature string
	Cert      []byte
	Report    []byte
}

// 去中心化的机密注入
type Secrets struct {
	Files map[string]string
	Env   map[string]string
}
