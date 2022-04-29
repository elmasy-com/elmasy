package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

// QueryTXT returns the answer as a string slice.
// Returns nil in case of error.
// This function retries the query in case of error, up to MAX_RETRIES.
func QueryTXT(name string) ([]string, error) {

	var (
		a   []dns.RR
		r   = make([]string, 0)
		err error
	)

	for i := 0; i < MAX_RETRIES; i++ {

		a, err = query(name, dns.TypeTXT)
		if err == nil {
			break
		}
	}

	for i := range a {

		switch v := a[i].(type) {
		case *dns.TXT:
			r = append(r, v.Txt...)
		default:
			return r, fmt.Errorf("unknown type: %T", v)
		}
	}

	return r, err

}
