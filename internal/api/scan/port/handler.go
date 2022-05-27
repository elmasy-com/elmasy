package port

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/elmasy/pkg/portscan"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	var (
		result portscan.Result
		errs   []error
	)
	params, err := parseQuery(c)
	if err != nil {
		return
	}

	switch params.Technique {
	case "stealth", "syn":
		result, errs = portscan.StealthScan(params.IP, []int{params.Port}, params.Timeout)
	case "connect":
		result, errs = portscan.ConnectScan(params.IP, []int{params.Port}, params.Timeout)
	}

	if len(errs) > 0 {

		for i := range errs {
			c.Error(errs[i])
		}

		// Return only the first error
		c.JSON(http.StatusInternalServerError, sdk.Error{Err: errs[0].Error()})
		return
	}

	if len(result) != 1 {
		c.Error(fmt.Errorf("Invalid number of result at single port: %d", len(result)))
		c.Error(fmt.Errorf("%#v", result))
	}

	c.JSON(http.StatusOK, sdk.ResultStr{Result: result[0].State.String()})
}
