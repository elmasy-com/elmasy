package utils

import (
	"fmt"
	"strings"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/identify"
)

// Lookup46 do a IPv4/IPv6 lookup on the domain.
// Return NXDOMAIN if both 4 and 6 is NXDOMAIN.
// This function adds brackets around IPv6 addresses.
func Lookup46(domain string) ([]string, error) {

	result := make([]string, 0)

	t, err := sdk.DNSLookup("A", domain)
	if err != nil && err.Error() != "NXDOMAIN" {
		return nil, err
	} else {
		result = append(result, t...)
	}

	t, err = sdk.DNSLookup("AAAA", domain)
	if err != nil && err.Error() != "NXDOMAIN" {
		return nil, err
	} else {
		result = append(result, t...)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("NXDOMAIN")
	}

	return result, nil
}

// IPv6BracketAdd add brackets around an IPv6 address.
// If ip is an IPv4 address, returns it as is.
func IPv6BracketAdd(ip string) string {

	// Add brackets to IPv6 addresses
	if identify.IsValidIPv6(ip) {
		return "[" + ip + "]"
	}

	return ip
}

// IPv6BracketRemove removes brackets around an IPv6 address.
// If ip is an IPv4 address, returns it as is.
func IPv6BracketRemove(ip string) string {

	return strings.Trim(ip, "[]")
}
