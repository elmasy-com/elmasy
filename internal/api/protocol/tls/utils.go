package tls

import (
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"
	"github.com/elmasy-com/identify"
)

// Choose a ServerName.
// The priority: servername parameter > target (if valid domain) > "" (empty string)
func getServerName(servername, target string) string {

	if servername != "" {
		return servername
	}

	if identify.IsDomainName(target) {
		return target
	}

	return ""
}

func resultCiphers[T ciphersuite.CipherSuite | sdk.Cipher](ciphers []T) []Cipher {

	r := make([]Cipher, 0)

	switch t := any(ciphers).(type) {
	case []ciphersuite.CipherSuite:
		for i := range t {
			r = append(r, Cipher{Name: t[i].Name, Security: t[i].Security})
		}
	case []sdk.Cipher:
		for i := range t {
			r = append(r, Cipher{Name: t[i].Name, Security: t[i].Security})
		}
	}

	return r
}
