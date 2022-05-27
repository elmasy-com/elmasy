package ip

import (
	"net/http"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, sdk.ResultStr{Result: c.ClientIP()})
	} else {
		c.String(http.StatusOK, c.ClientIP())
	}
}
