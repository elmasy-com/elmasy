package tls10

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"
)

/*
	struct {
	    ProtocolVersion server_version;
	    Random random;
	    SessionID session_id;
	    CipherSuite cipher_suite;
		CompressionMethod compression_method;
	} ServerHello
*/

type serverHello struct {
	Version           []byte
	Random            random
	SessionID         []byte
	CipherSuite       ciphersuite.CipherSuite
	CompressionMethod uint8
}

func unmarshalServerHello(bytes []byte) (serverHello, error) {

	var (
		hello serverHello
		ok    bool
		err   error
		buf   = bytebuilder.NewBuffer(bytes)
	)

	if hello.Version = buf.ReadBytes(2); hello.Version == nil {
		return hello, fmt.Errorf("failed to read Version")
	}
	if err := checkVersion(hello.Version[0], hello.Version[1]); err != nil {
		return hello, err
	}

	if hello.Random, err = unmarshalRandom(buf.ReadBytes(32)); err != nil {
		return hello, fmt.Errorf("failed to read Random: %s", err)
	}

	if hello.SessionID, ok = buf.ReadVector(8); !ok {
		return hello, fmt.Errorf("failed to read SessionID")
	}

	if hello.CipherSuite, err = ciphersuite.Unmarhsal(buf.ReadBytes(2)); err != nil {
		return hello, fmt.Errorf("failed to read CipherSuite: %s", err)
	}

	if hello.CompressionMethod, ok = buf.ReadUint8(); !ok {
		return hello, fmt.Errorf("failed to read CompressionMethod")
	}

	return hello, nil
}

func checkVersion(major, minor byte) error {

	if major == 0x54 && minor == 0x54 {
		return fmt.Errorf("unencrypted HTTP response")
	}

	if major != 0x03 {
		return fmt.Errorf("invalid protocol version: 0x%02x 0x%02x", major, minor)
	}

	if minor != 0x00 && minor != 0x01 && minor != 0x02 && minor != 0x03 {
		return fmt.Errorf("invalid protocol version: 0x%02x 0x%02x", major, minor)
	}

	return nil
}
