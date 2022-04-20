package ip

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"result": c.ClientIP()})
	} else {
		c.String(http.StatusOK, c.ClientIP())
	}
}
