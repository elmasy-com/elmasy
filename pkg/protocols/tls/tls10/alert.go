package tls10

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
)

/*
	enum { warning(1), fatal(2), (255) } AlertLevel;

	enum {
        close_notify(0),
        unexpected_message(10),
        bad_record_mac(20),
        decryption_failed(21),
        record_overflow(22),
        decompression_failure(30),
        handshake_failure(40),
        bad_certificate(42),
        unsupported_certificate(43),
        certificate_revoked(44),
        certificate_expired(45),
        certificate_unknown(46),
        illegal_parameter(47),
        unknown_ca(48),
        access_denied(49),
        decode_error(50),
        decrypt_error(51),
        export_restriction(60),
        protocol_version(70),
        insufficient_security(71),
        internal_error(80),
        user_canceled(90),
        no_renegotiation(100),
        (255)
    } AlertDescription;

	struct {
	    AlertLevel level;
	    AlertDescription description;
	} Alert;
*/

type alert struct {
	Level       uint8
	Description uint8
}

func marshalAlert(level, description uint8) []byte {

	return []byte{level, description}
}

func unmarshalAlert(bytes []byte) (alert, error) {

	var (
		buf = bytebuilder.NewBuffer(bytes)
		a   alert
		ok  bool
	)

	if a.Level, ok = buf.ReadUint8(); !ok {
		return a, fmt.Errorf("failed to read Level")
	}

	if a.Description, ok = buf.ReadUint8(); !ok {
		return a, fmt.Errorf("failed to read Description")
	}

	if buf.Size() != 0 {
		return a, fmt.Errorf("buf is not empty")
	}

	return a, nil
}
