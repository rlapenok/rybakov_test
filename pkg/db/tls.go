package db

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"

	pkgErr "github.com/rlapenok/rybakov_test/pkg/errors"
)

// errTLSConfigBuild - create new TLS config build error
func errTLSConfigBuild(err error) *pkgErr.Error {
	meta := map[string]any{
		"original_error": err.Error(),
	}
	return pkgErr.NewError(
		pkgErr.CodeInternalError,
		"Failed to build TLS config",
		"FAILED_TO_BUILD_TLS_CONFIG",
		meta,
	)
}

// BuildTLSConfig - build TLS config
func BuildTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errTLSConfigBuild(err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	if caPath != "" {
		caBytes, err := os.ReadFile(caPath)
		if err != nil {
			return nil, errTLSConfigBuild(err)
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caBytes) {
			return nil, errTLSConfigBuild(errors.New("invalid CA PEM"))
		}
		tlsConfig.RootCAs = caPool
		tlsConfig.InsecureSkipVerify = false
	}

	return tlsConfig, nil
}
