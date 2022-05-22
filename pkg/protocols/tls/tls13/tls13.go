package tls13

import (
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"
	tls "github.com/refraction-networking/utls"
)

type TLS13 struct {
	Supported     bool
	Certificates  []x509.Certificate
	DefaultCipher ciphersuite.CipherSuite
	Ciphers       []ciphersuite.CipherSuite
}

type Opts struct {
	ServerName string // SNI
}

func ciphersToUint16(ciphers []ciphersuite.CipherSuite) []uint16 {

	v := make([]uint16, 0)

	for i := range ciphers {

		v = append(v, binary.BigEndian.Uint16(ciphers[i].Value))
	}

	return v
}

func handshake(network, ip, port string, timeout time.Duration, ciphers []ciphersuite.CipherSuite, opts Opts) (TLS13, error) {

	var (
		conf   tls.Config
		result TLS13
	)

	if opts.ServerName == "" {
		conf.InsecureSkipVerify = true
	} else {
		conf.ServerName = opts.ServerName
	}

	dialConn, err := net.DialTimeout(network, ip+":"+port, timeout)
	if err != nil {
		return result, err
	}
	uTlsConn := tls.UClient(dialConn, &conf, tls.HelloCustom)
	defer uTlsConn.Close()

	spec := tls.ClientHelloSpec{
		TLSVersMax:   tls.VersionTLS13,
		TLSVersMin:   tls.VersionTLS10,
		CipherSuites: ciphersToUint16(ciphers),
		Extensions: []tls.TLSExtension{
			&tls.SupportedCurvesExtension{Curves: []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384, tls.CurveP521}},
			&tls.SupportedPointsExtension{SupportedPoints: []byte{0}}, // uncompressed
			&tls.SessionTicketExtension{},
			&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
				tls.ECDSAWithP256AndSHA256,
				tls.ECDSAWithP384AndSHA384,
				tls.ECDSAWithP521AndSHA512,
				tls.PSSWithSHA256,
				tls.PSSWithSHA384,
				tls.PSSWithSHA512,
				tls.PKCS1WithSHA256,
				tls.PKCS1WithSHA384,
				tls.PKCS1WithSHA512,
				tls.ECDSAWithSHA1,
				tls.PKCS1WithSHA1}},
			&tls.KeyShareExtension{KeyShares: []tls.KeyShare{{Group: tls.X25519}}},
			&tls.PSKKeyExchangeModesExtension{Modes: []uint8{1}}, // pskModeDHE
			&tls.SupportedVersionsExtension{Versions: []uint16{tls.VersionTLS13}},
		},
		GetSessionID: nil,
	}

	if opts.ServerName != "" {
		spec.Extensions = append(spec.Extensions, &tls.SNIExtension{})
	}

	if err := uTlsConn.ApplyPreset(&spec); err != nil {
		return result, fmt.Errorf("failed to apply spec: %s", err)
	}

	err = uTlsConn.Handshake()
	if err != nil {

		switch true {
		case strings.Contains(err.Error(), "handshake failure"):
			return result, nil
		case strings.Contains(err.Error(), "protocol version not supported"):
			return result, nil
		case strings.Contains(err.Error(), "EOF"):
			// Based on the tests, EOF means that no reaction to handshake, a "close notify" or TCP RST.
			return result, nil
		default:
			return result, fmt.Errorf("failed to handshake: %s", err)
		}
	}

	result.Supported = true

	for i := range uTlsConn.ConnectionState().PeerCertificates {
		result.Certificates = append(result.Certificates, *uTlsConn.ConnectionState().PeerCertificates[i])
	}

	if c := ciphersuite.FindByUint16(ciphers, uTlsConn.ConnectionState().CipherSuite); c == nil {
		return result, fmt.Errorf("failed to find ciphersuite: %x", uTlsConn.ConnectionState().CipherSuite)
	} else {
		result.DefaultCipher = *c
		result.Ciphers = append(result.Ciphers, *c)
	}

	return result, nil
}

// There are ciphersuites, that uTLS cant handle.
// In this case, iterate over it, one by one. If the error message is "server chose an unconfigured cipher suite", the ciphersuite is supported by the server.
func getUnconfiguredCiphers(network, ip, port string, timeout time.Duration, ciphers []ciphersuite.CipherSuite, opts Opts) ([]ciphersuite.CipherSuite, error) {

	supported := make([]ciphersuite.CipherSuite, 0)

	for i := range ciphers {

		_, err := handshake(network, ip, port, timeout, []ciphersuite.CipherSuite{ciphers[i]}, opts)

		if err != nil {
			if strings.Contains(err.Error(), "server chose an unconfigured cipher suite") {
				supported = append(supported, ciphers[i])
			} else {
				return supported, err
			}
		}
	}

	return supported, nil
}

func getSupportedCiphers(network, ip, port string, timeout time.Duration, ciphers []ciphersuite.CipherSuite, opts Opts) ([]ciphersuite.CipherSuite, error) {

	var (
		supported = make([]ciphersuite.CipherSuite, 0)
	)

	for {

		result, err := handshake(network, ip, port, timeout, ciphers, opts)
		if err != nil {

			if strings.Contains(err.Error(), "server chose an unconfigured cipher suite") {
				unconfigured, err := getUnconfiguredCiphers(network, ip, port, timeout, ciphers, opts)
				if err != nil {
					return supported, fmt.Errorf("failed to get unconfigured ciphers: %s", err)
				}
				supported = append(supported, unconfigured...)
				return supported, nil
			}

			return supported, fmt.Errorf("failed to do handshake: %s", err)
		}

		if !result.Supported {
			return supported, nil
		}

		ciphers = ciphersuite.Remove(ciphers, result.DefaultCipher)
		supported = append(supported, result.DefaultCipher)
	}
}

func Scan(network, ip, port string, timeout time.Duration, opts Opts) (TLS13, error) {

	ciphers := ciphersuite.Get(ciphersuite.TLS13)

	result, err := handshake(network, ip, port, timeout, ciphers, opts)
	if err != nil {
		return result, fmt.Errorf("handshake failed: %s", err)
	}

	if !result.Supported {
		return result, nil
	}

	ciphers = ciphersuite.Remove(ciphers, result.DefaultCipher)

	supported, err := getSupportedCiphers(network, ip, port, timeout, ciphers, opts)
	if err != nil {
		return result, fmt.Errorf("failed to get supported ciphers: %s", err)
	}

	result.Ciphers = append(result.Ciphers, supported...)

	return result, nil

}

func Probe(network, ip, port string, timeout time.Duration, opts Opts) (bool, error) {

	ciphers := ciphersuite.Get(ciphersuite.TLS13)

	r, err := handshake(network, ip, port, timeout, ciphers, opts)

	return r.Supported, err
}
