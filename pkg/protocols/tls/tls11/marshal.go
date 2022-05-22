package tls11

import "github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"

// Shorthand to create a Closure Alert
func createClosureAlert() []byte {

	fragment := marshalAlert(1, 0)

	return marshalSSLPlaintext(21, fragment)
}

// Shorthand to create a ClientHello
func createPacketClientHello(ciphers []ciphersuite.CipherSuite, ServerName string) []byte {

	body := marshalClientHello(ciphers, ServerName)

	fragment := marshalHandshake(1, body)

	return marshalSSLPlaintext(22, fragment)
}
