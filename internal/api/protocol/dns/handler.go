package dns

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/elmasy/pkg/protocols/dns"
	"github.com/elmasy-com/identify"
	"github.com/elmasy-com/slices"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	name := c.Param("name")
	if !identify.IsDomainName(name) {
		err := fmt.Errorf("Invalid name: %s", name)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	switch qtype := c.Param("type"); qtype {
	case "A":
		r, err := dns.QueryA(name)
		if err != nil {
			c.Error(err)
			c.JSON(getStatusCode(err), sdk.Error{Err: err.Error()})
			return
		}

		c.JSON(http.StatusOK, sdk.ResultStrs{Results: slices.Strings(r)})

	case "AAAA":
		r, err := dns.QueryAAAA(name)
		if err != nil {
			c.Error(err)
			c.JSON(getStatusCode(err), sdk.Error{Err: err.Error()})
			return
		}

		c.JSON(http.StatusOK, sdk.ResultStrs{Results: slices.Strings(r)})

	case "MX":
		r, err := dns.QueryMX(name)
		if err != nil {
			c.Error(err)
			c.JSON(getStatusCode(err), sdk.Error{Err: err.Error()})
			return
		}

		c.JSON(http.StatusOK, sdk.ResultStrs{Results: mxToString(r)})

	case "TXT":
		r, err := dns.QueryTXT(name)
		if err != nil {
			c.Error(err)
			c.JSON(getStatusCode(err), sdk.Error{Err: err.Error()})
			return
		}

		c.JSON(http.StatusOK, sdk.ResultStrs{Results: r})

	default:
		err := fmt.Errorf("Invalid type: %s", qtype)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
	}
}
