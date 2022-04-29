package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

type MX struct {
	Preference int    // Priority
	Exchange   string // Server's hostname
}

// QueryMX returns MX structs of the answers.
// Returns nil in case of error.
// This function retries the query in case of error, up to MAX_RETRIES.
func QueryMX(name string) ([]MX, error) {

	var (
		a   []dns.RR
		r   = make([]MX, 0)
		err error
	)

	for i := 0; i < MAX_RETRIES; i++ {

		a, err = query(name, dns.TypeMX)
		if err == nil {
			break
		}
	}

	for i := range a {

		switch v := a[i].(type) {
		case *dns.MX:
			r = append(r, MX{Preference: int(v.Preference), Exchange: v.Mx})
		default:
			return r, fmt.Errorf("unknown type: %T", v)
		}
	}

	return r, err

}
