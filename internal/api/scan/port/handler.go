package port

import (
	"fmt"
	"net/http"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs[0]})
		return
	}

	if len(result) != 1 {
		c.Error(fmt.Errorf("Multiple result at single port: %#v", result))
	}

	c.JSON(http.StatusOK, gin.H{"result": result[0].State.String()})
}
