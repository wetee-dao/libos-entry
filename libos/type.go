package libos

type LoadParam struct {
	Address   string
	Time      string
	Signature string
}

type Secrets struct {
	Files map[string]string
	Env   map[string]string
}
