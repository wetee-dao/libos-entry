package ego

// const (
// 	// LibosEnvironmentRootCA 包含持有PEM编码根证书的环境变量名称
// 	LibosEnvironmentRootCA = "LIBOS_ROOT_CA"
// 	// LibosEnvironmentPrivateKey 包含持有属于 libos 特定证书的PEM编码私钥的环境变量名称
// 	LibosEnvironmentPrivateKey = "LIBOS_PRIVATE_KEY"
// 	// LibosEnvironmentCertificateChain 包含持有 libos 特定 PEM 编码证书的环境变量名称
// 	LibosEnvironmentCertificateChain = "LIBOS_CERTIFICATE_CHAIN"
// )

// // 为libos提供了一个预配置的TLS配置,使用本地 worker 作为信任锚点
// func GetTLSConfig(verifyClientCerts bool) (*tls.Config, error) {
// 	tlsCert, roots, err := getRootEnv()
// 	if err != nil {
// 		return nil, err
// 	}

// 	tlsConfig := &tls.Config{
// 		RootCAs:      roots,
// 		Certificates: []tls.Certificate{tlsCert},
// 	}

// 	if verifyClientCerts {
// 		tlsConfig.ClientCAs = roots
// 		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
// 	}

// 	return tlsConfig, nil
// }

// // 从容器环境变量中获取TLS证书
// func getRootEnv() (tls.Certificate, *x509.CertPool, error) {
// 	certChain, err := getByteEnv(LibosEnvironmentCertificateChain)
// 	if err != nil {
// 		return tls.Certificate{}, nil, err
// 	}
// 	libosRootCA, err := getByteEnv(LibosEnvironmentRootCA)
// 	if err != nil {
// 		return tls.Certificate{}, nil, err
// 	}
// 	leafPrivk, err := getByteEnv(LibosEnvironmentPrivateKey)
// 	if err != nil {
// 		return tls.Certificate{}, nil, err
// 	}

// 	roots := x509.NewCertPool()
// 	if !roots.AppendCertsFromPEM(libosRootCA) {
// 		return tls.Certificate{}, nil, fmt.Errorf("cannot append libosRootCA to CertPool")
// 	}

// 	tlsCert, err := tls.X509KeyPair(certChain, leafPrivk)
// 	if err != nil {
// 		return tls.Certificate{}, nil, fmt.Errorf("cannot create TLS cert: %v", err)
// 	}

// 	return tlsCert, roots, nil
// }

// // 获取环境变量的值
// func getByteEnv(name string) ([]byte, error) {
// 	value := os.Getenv(name)
// 	if len(value) == 0 {
// 		return nil, fmt.Errorf("environment variable not set: %s", name)
// 	}
// 	return []byte(value), nil
// }
