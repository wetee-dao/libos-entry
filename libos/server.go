package libos

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"net/http"
)

func startEntryServer(cert []byte, priv crypto.PrivateKey) error {
	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}

	server := http.Server{Addr: "0.0.0.0:8883", TLSConfig: &tlsCfg}
	fmt.Println("Start entry secret listening 0.0.0.0:8883 ...")
	err := server.ListenAndServeTLS("", "")
	return err
}
