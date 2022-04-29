package ssl30

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
	"github.com/elmasy-com/protocols/tls/ciphersuite"
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
	} else if hello.Version[0] != VER_MAJOR || hello.Version[1] != VER_MINOR {
		return hello, fmt.Errorf("invalid server version: 0x%02x 0x%02x", hello.Version[0], hello.Version[1])
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
