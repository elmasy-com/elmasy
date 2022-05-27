package tls

import (
	"crypto/x509"
	"fmt"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ssl30"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls10"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls11"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls12"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls13"
)

type TLS struct {
	Supported     bool
	Certificates  []x509.Certificate
	DefaultCipher ciphersuite.CipherSuite
	Ciphers       []ciphersuite.CipherSuite
}

func Scan(version, network, ip, port string, timeout time.Duration, servername string) (TLS, error) {
	switch version {
	case "ssl30":
		r, err := ssl30.Scan(network, ip, port, timeout)
		return TLS(r), err
	case "tls10":
		r, err := tls10.Scan(network, ip, port, timeout, servername)
		return TLS(r), err
	case "tls11":
		r, err := tls11.Scan(network, ip, port, timeout, servername)
		return TLS(r), err
	case "tls12":
		r, err := tls12.Scan(network, ip, port, timeout, servername)
		return TLS(r), err
	case "tls13":
		r, err := tls13.Scan(network, ip, port, timeout, servername)
		return TLS(r), err
	default:
		return TLS{}, fmt.Errorf("invalid version: %s", version)
	}
}

func Handshake(version, network, ip, port string, timeout time.Duration, servername string) (TLS, error) {

	switch version {
	case "ssl30":
		r, err := ssl30.Handshake(network, ip, port, timeout)
		return TLS(r), err
	case "tls10":
		r, err := tls10.Handshake(network, ip, port, timeout, servername)
		return TLS(r), err
	case "tls11":
		r, err := tls11.Handshake(network, ip, port, timeout, servername)
		return TLS(r), err
	case "tls12":
		r, err := tls12.Handshake(network, ip, port, timeout, servername)
		return TLS(r), err
	case "tls13":
		r, err := tls13.Handshake(network, ip, port, timeout, servername)
		return TLS(r), err
	default:
		return TLS{}, fmt.Errorf("invalid version: %s", version)
	}
}

func Probe(version, network, ip, port string, timeout time.Duration, servername string) (bool, error) {

	switch version {
	case "ssl30":
		r, err := ssl30.Probe(network, ip, port, timeout)
		return r, err
	case "tls10":
		r, err := tls10.Probe(network, ip, port, timeout, servername)
		return r, err
	case "tls11":
		r, err := tls11.Probe(network, ip, port, timeout, servername)
		return r, err
	case "tls12":
		r, err := tls12.Probe(network, ip, port, timeout, servername)
		return r, err
	case "tls13":
		r, err := tls13.Probe(network, ip, port, timeout, servername)
		return r, err
	default:
		return false, fmt.Errorf("invalid version: %s", version)
	}
}
