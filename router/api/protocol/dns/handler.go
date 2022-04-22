package dns

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/dns"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	name := c.Param("name")
	if !identify.IsDomainName(name) {
		err := fmt.Errorf("Invalid name: %s", name)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch qtype := c.Param("type"); qtype {
	case "A":
		r, err := dns.QueryA(name)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": r})
	case "AAAA":
		r, err := dns.QueryAAAA(name)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": r})
	case "MX":
		r, err := dns.QueryMX(name)
		if err != nil {
			handleError(c, err)

			return
		}
		c.JSON(http.StatusOK, gin.H{"results": mxToString(r)})
	case "TXT":
		r, err := dns.QueryTXT(name)
		if err != nil {
			handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"results": r})
	default:
		err := fmt.Errorf("Invalid type: %s", qtype)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Error(err)
	}
}
