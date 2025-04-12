package tools

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"time"
)

func LoadClientPool(cert []byte) (*x509.CertPool, error) {
	// Create a certificate pool and add the server's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, errors.New("failed to add server CA's certificate")
	}
	return certPool, nil
}

func CreateHTTPClient(certPool *x509.CertPool) *http.Client {
	tlsConfig := &tls.Config{RootCAs: certPool}
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
