package tools

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

func CreateTlsConfig(clientCert, clientKey, serverCA []byte) (*tls.Config, error) {
	clientCertPair, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse client certificate and key: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(serverCA) {
		return nil, fmt.Errorf("failed to add server CA certificate to pool")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCertPair},
		RootCAs:      caCertPool,
	}
	return config, nil
}

func CreateHTTPClient(tlsConfig *tls.Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Protocols:       new(http.Protocols),
	}
	transport.Protocols.SetHTTP2(true)

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}
