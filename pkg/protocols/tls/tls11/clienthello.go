package tls11

import (
	"github.com/elmasy-com/bytebuilder"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"
)

/*
	struct {
	    uint32 gmt_unix_time;
	    opaque random_bytes[28];
	} Random;

	struct {
		ProtocolVersion client_version;
		Random random;
		SessionID session_id;
		CipherSuite cipher_suites<2..2^16-1>;
		CompressionMethod compression_methods<1..2^8-1>;
	} ClientHello;
*/

func marshalClientHello(ciphers []ciphersuite.CipherSuite, ServerName string) []byte {

	buf := bytebuilder.NewEmpty()

	buf.WriteBytes(VER_MAJOR, VER_MINOR)

	buf.WriteBytes(marshalRandom()...)

	buf.WriteBytes(marshalSessionID()...)

	buf.WriteVector(ciphersuite.Marshal(ciphers), 16)

	buf.WriteVector([]byte{0}, 8)

	exts := marshalExtensions(ServerName)

	// Extensions length
	buf.WriteUint16(uint16(len(exts)))
	buf.WriteBytes(exts...)

	return buf.Bytes()
}
