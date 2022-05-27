package dns

import (
	"net/http"
	"sort"
	"strings"

	"github.com/elmasy-com/elmasy/pkg/protocols/dns"
)

// mxToString sort the MX records by preference (or by name if preferences are equal)
// and returns the names (Exchange) as a string slice
func mxToString(mxs []dns.MX) []string {

	r := make([]string, 0)

	// If preference is equal, compare the name
	sort.Slice(mxs, func(i, j int) bool {
		if mxs[i].Preference == mxs[j].Preference {
			return mxs[i].Exchange < mxs[j].Exchange
		}
		return mxs[i].Preference < mxs[j].Preference
	})

	for i := range mxs {
		// MX returns fqdn so trim the trailing "."
		r = append(r, strings.TrimSuffix(mxs[i].Exchange, "."))
	}

	return r
}

// Returns the status code accoding to the error message
func getStatusCode(err error) int {

	switch err.Error() {
	case "FORMERR", "SERVFAIL", "NOTIMP", "REFUSED", "YXDOMAIN", "XRRSET", "NOTAUTH", "NOTZONE":
		return http.StatusInternalServerError
	case "NXDOMAIN":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError

	}
}
