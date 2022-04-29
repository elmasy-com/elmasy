package tls11

import "github.com/elmasy-com/protocols/tls/ciphersuite"

// Shorthand to create a Closure Alert
func createClosureAlert() []byte {

	fragment := marshalAlert(1, 0)

	return marshalSSLPlaintext(21, fragment)
}

// Shorthand to create a ClientHello
func createPacketClientHello(ciphers []ciphersuite.CipherSuite) []byte {

	body := marshalClientHello(ciphers)

	fragment := marshalHandshake(1, body)

	return marshalSSLPlaintext(22, fragment)
}
