package ssl30

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
)

/*
	enum {
	    hello_request(0), client_hello(1), server_hello(2),
	    certificate(11), server_key_exchange (12),
	    certificate_request(13), server_hello_done(14),
	    certificate_verify(15), client_key_exchange(16),
	    finished(20), (255)
	} HandshakeType;

	struct {
		HandshakeType msg_type;    // handshake type
	    uint24 length;             // bytes in message
	    select (HandshakeType) {
	        case hello_request: HelloRequest;
	        case client_hello: ClientHello;
	        case server_hello: ServerHello;
	        case certificate: Certificate;
	        case server_key_exchange: ServerKeyExchange;
	        case certificate_request: CertificateRequest;
	        case server_hello_done: ServerHelloDone;
	        case certificate_verify: CertificateVerify;
	        case client_key_exchange: ClientKeyExchange;
	        case finished: Finished;
	    } body;
	} Handshake;
*/

func marshalHandshake(msgType uint8, body []byte) []byte {

	buf := bytebuilder.NewEmpty()

	buf.WriteUint8(msgType)
	buf.WriteInt(len(body), 24)
	buf.WriteBytes(body...)

	return buf.Bytes()
}

// Because SSLPlaintext can contain multiple handshake, this function returns a slice.
// Iterates over and over on bytes, until every handshake is read or any error occur.
func unmarshalHandshake(bytes []byte) ([]interface{}, error) {

	var (
		buf      = bytebuilder.NewBuffer(bytes)
		msgType  uint8
		length   uint32
		ok       bool
		body     []byte
		messages []interface{}
		err      error
	)

	for !buf.Empty() {

		var message interface{}

		if msgType, ok = buf.ReadUint8(); !ok {
			return messages, fmt.Errorf("failed to read msgType")
		}

		if length, ok = buf.ReadUint24(); !ok {
			return messages, fmt.Errorf("failed to read length")
		}

		if body = buf.ReadBytes(int(length)); body == nil {
			return messages, fmt.Errorf("failed to read body")
		}

		switch msgType {
		case 0:
			return messages, fmt.Errorf("handshake type hello_request is not supported")
		case 1:
			return messages, fmt.Errorf("handshake type client_hello is invalid")
		case 2:
			if message, err = unmarshalServerHello(body); err != nil {
				return messages, fmt.Errorf("failed to unmarshal ServerHello: %s", err)
			}

		case 11:
			if message, err = unmarhsalCertificate(body); err != nil {
				return messages, fmt.Errorf("failed to unmarshal Certificate: %s", err)
			}
		case 12:
			if message, err = unmarshalServerKeyExchange(body); err != nil {
				return messages, fmt.Errorf("failed to unmarshal ServerKeyExchange: %s", err)
			}
		case 13:
			if message, err = unmarshalCertificateRequest(body); err != nil {
				return messages, fmt.Errorf("failed to unmarshal CertificateRequest: %s", err)
			}
		case 14:
			message = serverHelloDone{}
		case 15:
			return messages, fmt.Errorf("handshake type certificate_verify is not supported")
		case 16:
			return messages, fmt.Errorf("handshake type client_key_exchange is not supported")
		case 20:
			return messages, fmt.Errorf("handshake type finished is not supported")
		default:
			return messages, fmt.Errorf("unknown Handshake type: %d", msgType)
		}

		messages = append(messages, message)

	}

	return messages, nil
}
