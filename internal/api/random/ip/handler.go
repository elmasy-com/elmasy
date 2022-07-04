package randomip

import (
	"net/http"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/elmasy/pkg/randomip"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	switch version := c.Query("version"); version {
	case "ipv4", "4":
		c.JSON(http.StatusOK, sdk.ResultStr{Result: randomip.GetPublicIPv4().String()})
	case "ipv6", "6":
		c.JSON(http.StatusOK, sdk.ResultStr{Result: randomip.GetPublicIPv6().String()})
	default:
		c.JSON(http.StatusOK, sdk.ResultStr{Result: randomip.GetPublicIP().String()})
	}
}
