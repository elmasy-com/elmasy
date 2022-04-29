package tls11

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
)

/*
	struct {
		uint8 major, minor;
	} ProtocolVersion;

	enum {
	    change_cipher_spec(20), alert(21), handshake(22),
	    application_data(23), (255)
	} ContentType;

	struct {
	    ContentType type;
	    ProtocolVersion version;
	    uint16 length;
	    opaque fragment[SSLPlaintext.length];
	} SSLPlaintext;
*/

func marshalSSLPlaintext(cType uint8, fragment []byte) []byte {

	buf := bytebuilder.NewEmpty()

	buf.WriteUint8(cType)
	buf.WriteBytes(VER_MAJOR, VER_MINOR)
	buf.WriteInt(len(fragment), 16)
	buf.WriteBytes(fragment...)

	return buf.Bytes()
}

func unmarshalSSLPlaintext(buf *bytebuilder.Buffer) ([]interface{}, error) {

	var (
		err        error
		cType      uint8
		length     uint16
		ok         bool
		fragment   []byte
		handshakes []interface{}
	)

	if cType, ok = buf.ReadUint8(); !ok {
		return handshakes, fmt.Errorf("failed to read Length")
	}

	if version := buf.ReadBytes(2); version == nil {
		return handshakes, fmt.Errorf("failed to read protocol version")
	} else if version[0] != VER_MAJOR || version[1] != VER_MINOR {
		return handshakes, fmt.Errorf("invalid protocol version: 0x%02X 0x%02X", version[0], version[1])
	}

	if length, ok = buf.ReadUint16(); !ok {
		return handshakes, fmt.Errorf("failed to read length")
	}

	if fragment = buf.ReadBytes(int(length)); fragment == nil {
		return handshakes, fmt.Errorf("failed to read fragment")
	}

	switch cType {
	case 20:
		return handshakes, fmt.Errorf("SSLPlaintext type change_cipher_spec is not supported")
	case 21:
		var handshake interface{}
		if handshake, err = unmarshalAlert(fragment); err != nil {
			return handshakes, fmt.Errorf("failed to unmarshal Alert: %s", err)
		}
		handshakes = append(handshakes, handshake)
	case 22:
		if handshakes, err = unmarshalHandshake(fragment); err != nil {
			return handshakes, fmt.Errorf("failed to unmarshal SSLPlaintext: %s", err)
		}
	case 23:
		return handshakes, fmt.Errorf("SSLPlaintext type application_data is not supported")
	default:
		return handshakes, fmt.Errorf("SSLPlaintext type unknown: %d", cType)
	}

	return handshakes, nil
}
