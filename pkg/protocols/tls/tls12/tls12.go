package tls12

import (
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/elmasy-com/bytebuilder"
	"github.com/elmasy-com/protocols/tls/ciphersuite"
)

const (
	VER_MAJOR uint8 = 0x03
	VER_MINOR uint8 = 0x03
)

type TLS12 struct {
	Supported     bool
	Certificates  []x509.Certificate
	DefaultCipher ciphersuite.CipherSuite
	Ciphers       []ciphersuite.CipherSuite
}

func sendClientHello(conn *net.Conn, timeout time.Duration, ciphers []ciphersuite.CipherSuite) error {

	hello := createPacketClientHello(ciphers)

	if err := (*conn).SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %s", err)
	}

	if num, err := (*conn).Write(hello); err != nil {
		return fmt.Errorf("failed to write: %s", err)
	} else if num != len(hello) {
		return fmt.Errorf("fewer bytes writen: want(%d) / actual(%d)", len(hello), num)
	}

	return nil
}

func readServerResponse(conn *net.Conn, timeout time.Duration) ([]byte, error) {

	var (
		buf bytebuilder.Buffer
		err error
	)

	if err := (*conn).SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return []byte{}, fmt.Errorf("failed to set read deadline: %s", err)
	}

	if buf, err = bytebuilder.ReadAll(*conn); err != nil {

		if strings.Contains(err.Error(), "i/o timeout") {
			// Unresponsive server
			err = nil
		} else if strings.Contains(err.Error(), "connection reset by peer") && buf.Size() == 7 {
			// Some servers send an RST straight after a Alert(Handshake failure) packet *at the first handshake*.
			// The alert size (including SSLPLaintext) should be 7 byte.
			err = nil
		} else {
			err = fmt.Errorf("failed to read response: %s", err)
		}
	}

	return buf.Bytes(), err
}

func sendClosureALert(conn *net.Conn, timeout time.Duration) error {

	close := createClosureAlert()

	if err := (*conn).SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return fmt.Errorf("failed to set write deadline: %s", err)
	}

	if num, err := (*conn).Write(close); err != nil {
		return fmt.Errorf("failed to write: %s", err)
	} else if num != len(close) {
		return fmt.Errorf("fewer bytes writen: want(%d) / actual(%d)", len(close), num)
	}

	return nil
}

// Do the handshake and return the response as a byte slice.
func handshake(network, address string, timeout time.Duration, ciphers []ciphersuite.CipherSuite) (TLS12, error) {

	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return TLS12{}, fmt.Errorf("failed to connect to %s: %s", address, err)
	}
	defer conn.Close()

	if err := sendClientHello(&conn, timeout, ciphers); err != nil {
		return TLS12{}, fmt.Errorf("failed to send ClientHello: %s", err)
	}

	resp, err := readServerResponse(&conn, timeout)
	if err != nil {
		return TLS12{}, fmt.Errorf("failed to read server response: %s", err)
	}

	result, err := unmarshalResponse(resp)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal response: %s", err)
	}

	if result.Supported {
		if err := sendClosureALert(&conn, timeout); err != nil {
			return result, fmt.Errorf("failed to send Closure Alert: %s", err)
		}
	}

	return result, nil
}

func getSupportedCiphers(network, address string, timeout time.Duration, ciphers []ciphersuite.CipherSuite) ([]ciphersuite.CipherSuite, error) {

	var (
		supported = make([]ciphersuite.CipherSuite, 0)
	)

	for {

		result, err := handshake(network, address, timeout, ciphers)
		if err != nil && !strings.Contains(err.Error(), "connection reset by peer") {
			return supported, fmt.Errorf("failed to do handshake: %s", err)
		}

		if !result.Supported {
			return supported, nil
		}

		ciphers = ciphersuite.Remove(ciphers, result.DefaultCipher)
		supported = append(supported, result.DefaultCipher)
	}

}

func Scan(network, address string, timeout time.Duration) (TLS12, error) {

	ciphers := ciphersuite.Get(ciphersuite.TLS12)

	result, err := handshake(network, address, timeout, ciphers)
	if err != nil {
		return result, fmt.Errorf("handshake failed: %s", err)
	}

	if !result.Supported {
		return result, nil
	}

	ciphers = ciphersuite.Remove(ciphers, result.DefaultCipher)

	supported, err := getSupportedCiphers(network, address, timeout, ciphers)
	if err != nil {
		return result, fmt.Errorf("failed to get supported ciphers: %s", err)
	}

	result.Ciphers = append(result.Ciphers, result.DefaultCipher)
	result.Ciphers = append(result.Ciphers, supported...)

	return result, nil
}

func Probe(network, ip, port string, timeout time.Duration) (bool, error) {

	ciphers := ciphersuite.Get(ciphersuite.TLS12)

	r, err := handshake(network, ip+":"+port, timeout, ciphers)

	return r.Supported, err
}
