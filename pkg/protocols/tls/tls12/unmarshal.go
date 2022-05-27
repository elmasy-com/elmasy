package tls12

import (
	"github.com/elmasy-com/bytebuilder"
)

func unmarshalResponse(bytes []byte) (TLS12, error) {

	var (
		buf      = bytebuilder.NewBuffer(bytes)
		result   = TLS12{}
		messages []interface{}
	)

	for !buf.Empty() {

		message, err := unmarshalSSLPlaintext(&buf)
		if err != nil {
			return result, err
		}

		messages = append(messages, message...)

	}

	for i := range messages {

		switch message := messages[i].(type) {
		case alert:
			result.Supported = false
			return result, nil
		case serverHello:
			result.Supported = true
			result.DefaultCipher = message.CipherSuite
		case certificate:
			result.Certificates = message.Certificates

		}
	}

	return result, nil
}
