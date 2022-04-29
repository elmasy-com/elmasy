package ssl30

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
)

func unmarshalResponse(bytes []byte) (SSL30, error) {

	var (
		buf      = bytebuilder.NewBuffer(bytes)
		result   = SSL30{}
		messages []interface{}
	)

	for !buf.Empty() {

		message, err := unmarshalSSLPlaintext(&buf)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal SSLPlaintext: %s", err)
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
