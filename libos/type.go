package libos

type Secrets struct {
	Files map[string]string
	Env   map[string]string
}
