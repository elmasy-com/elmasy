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

// convertCiphers converts []ciphersuite.CipherSuite to []sdk.Cipher
func convertCiphers(ciphers []ciphersuite.CipherSuite) []sdk.Cipher {

	r := make([]sdk.Cipher, 0)

	for i := range ciphers {
		r = append(r, sdk.Cipher{Name: ciphers[i].Name, Security: ciphers[i].Security})
	}

	return r
}
