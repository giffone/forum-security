package config

import (
	"crypto/tls"
	"log"
)

const (
	ForumCertFile = "secure/cert/localhost.crt"
	ForumKeyFile  = "secure/cert/localhost.key"
)

func NewCrtConf() *tls.Config {
	cert, err := tls.LoadX509KeyPair(ForumCertFile, ForumKeyFile)
	if err != nil {
		log.Fatalf("server: certificate: load x509: %s", err.Error())
	}
	return &tls.Config{
		Certificates:             []tls.Certificate{cert},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}
}
