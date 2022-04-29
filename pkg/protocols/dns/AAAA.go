package dns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

// QueryAAAA returns net.IP structs of the answers.
// Returns nil in case of error.
// This function retries the query in case of error, up to MAX_RETRIES.
func QueryAAAA(name string) ([]net.IP, error) {

	var (
		a   []dns.RR
		r   = make([]net.IP, 0)
		err error
	)

	for i := 0; i < MAX_RETRIES; i++ {

		a, err = query(name, dns.TypeAAAA)
		if err == nil {
			break
		}
	}

	for i := range a {

		switch v := a[i].(type) {
		case *dns.AAAA:
			r = append(r, v.AAAA)
		case *dns.CNAME:
			// Ignore CNAME
			continue
		default:
			return r, fmt.Errorf("unknown type: %T", v)
		}
	}

	return r, err
}
