package port

import (
	"net/http"
	"time"

	"github.com/elmasy-com/portscan"
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
		result, errs = portscan.StealthScan(params.IP, params.Ports, 8*time.Second, 1*time.Second)
	case "connect":
		result, errs = portscan.ConnectScan(params.IP, params.Ports, 2*time.Second)
	}

	if len(errs) > 0 {
		errsStr := make([]string, 0)
		for i := range errs {
			c.Error(errs[i])
			errsStr = append(errsStr, errs[i].Error())
		}
		c.JSON(http.StatusInternalServerError, gin.H{"errors": errsStr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": convertResultString(result)})
}
