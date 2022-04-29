package ciphersuite

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
	"github.com/elmasy-com/slices"
)

const (
	SSL30 uint16 = 0x0300
	TLS10 uint16 = 0x0301
	TLS11 uint16 = 0x0302
	TLS12 uint16 = 0x0303
	TLS13 uint16 = 0x0304
)

type CipherSuite struct {
	Value   []byte
	Name    string
	version []uint16
}

// Get returns the known ciphers for version, or nil if the version is invalid.
func Get(version uint16) []CipherSuite {

	var v []CipherSuite

	for i := range CipherSuites {

		if slices.Contain(CipherSuites[i].version, version) {
			v = append(v, CipherSuites[i])
		}
	}

	return v
}

// Compare returns returns -1 if a < b, 0 if a == b or 1 if a > b.
func (a CipherSuite) Compare(b CipherSuite) int {

	switch {
	case a.Value[0] < b.Value[0]:
		return -1
	case a.Value[0] > b.Value[0]:
		return 1
	case a.Value[1] < b.Value[1]:
		return -1
	case a.Value[1] > b.Value[1]:
		return 1
	default:
		return 0
	}
}

func (c CipherSuite) String() string {
	return c.Name
}

// Marshal marshals ciphers to a byte slice
func Marshal(ciphers []CipherSuite) []byte {

	buf := bytebuilder.NewEmpty()

	for i := range ciphers {
		buf.WriteBytes(ciphers[i].Value...)
	}

	return buf.Bytes()
}

// Unmarshal returns a known CipherSuite from CipherSuites based on bytes.
func Unmarhsal(bytes []byte) (CipherSuite, error) {

	if bytes == nil {
		return CipherSuite{}, fmt.Errorf("bytes is nil")
	}
	if len(bytes) != 2 {
		return CipherSuite{}, fmt.Errorf("invalid length of bytes: %d", len(bytes))
	}

	for i := range CipherSuites {
		if CipherSuites[i].Value[0] == bytes[0] &&
			CipherSuites[i].Value[1] == bytes[1] {
			return CipherSuites[i], nil
		}
	}

	return CipherSuite{}, fmt.Errorf("unknown bytes: 0x%2.X 0x%2.X", bytes[0], bytes[1])
}

func Remove(ciphers []CipherSuite, cipher CipherSuite) []CipherSuite {

	v := make([]CipherSuite, 0)

	for i := range ciphers {
		if ciphers[i].Compare(cipher) != 0 {
			v = append(v, ciphers[i])
		}
	}

	return v
}
