package tls

import "github.com/elmasy-com/identify"

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
