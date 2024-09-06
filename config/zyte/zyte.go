package zyte

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Config struct {
	URL      string `env:"ZYTE_URL"`
	CertPath string `env:"ZYTE_CERT_PATH"`
}

func (cfg *Config) Cert() ([]byte, error) {
	certPath := cfg.CertPath
	if certPath == "" {
		path, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		certPath = filepath.Join(path, "zyte-proxy-ca.crt")
	}
	return os.ReadFile(certPath)
}

func Client(cfg Config) (*http.Client, error) {
	caCert, err := cfg.Cert()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	zyteproxyURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(zyteproxyURL),
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}
	return client, nil
}
