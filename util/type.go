package util

// 请求开始TEE容器
type TeeParam struct {
	Address string
	Time    int64
	Data    []byte
	Report  []byte
}

// 去中心化的机密注入
type Secrets struct {
	Files map[string]string
	Env   map[string]string
}
