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
		Certificates:       []tls.Certificate{clientCertPair},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			opts := x509.VerifyOptions{
				Roots: caCertPool,
			}

			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return fmt.Errorf("failed to parse certificate: %v", err)
			}

			_, err = cert.Verify(opts)
			if err != nil {
				return fmt.Errorf("certificate verification failed: %v", err)
			}

			return nil
		},
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
