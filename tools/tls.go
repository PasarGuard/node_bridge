package tools

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"
)

func LoadClientPool() (*x509.CertPool, error) {
	// Create a certificate pool and add the server's certificate
	certPool, err := x509.SystemCertPool()
	if err != nil {
		certPool = x509.NewCertPool()
	}
	//if !certPool.AppendCertsFromPEM(cert) {
	//	return nil, errors.New("failed to add server CA's certificate")
	//}
	return certPool, nil
}

func CreateHTTPClient(certPool *x509.CertPool, hostname string) *http.Client {
	tlsConfig := &tls.Config{RootCAs: certPool, ServerName: hostname}
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
