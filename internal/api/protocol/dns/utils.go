package dns

import (
	"net/http"
	"sort"

	"github.com/elmasy-com/elmasy/pkg/protocols/dns"
	"github.com/gin-gonic/gin"
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
		r = append(r, mxs[i].Exchange)
	}

	return r
}

func handleError(c *gin.Context, err error) {

	var code int

	switch err.Error() {
	case "FORMERR", "SERVFAIL", "NOTIMP", "REFUSED", "YXDOMAIN", "XRRSET", "NOTAUTH", "NOTZONE":
		code = http.StatusInternalServerError
	case "NXDOMAIN":
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError

	}

	c.Error(err)
	c.JSON(code, gin.H{"error": err.Error()})

}
