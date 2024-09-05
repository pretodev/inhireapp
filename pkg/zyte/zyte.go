package zyte

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func SetProxy(client *http.Client, proxyURL string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	certPath := filepath.Join(path, "zyte-proxy-ca.crt")
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	zyteproxyURL, err := url.Parse(proxyURL)
	if err != nil {
		return err
	}
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(zyteproxyURL),
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}
	return nil
}
