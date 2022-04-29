package randomip

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/elmasy/pkg/randomip"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	switch version := c.Param("version"); version {
	case "ipv4", "4":
		c.JSON(http.StatusOK, gin.H{"result": randomip.GetPublicIPv4()})
	case "ipv6", "6":
		c.JSON(http.StatusOK, gin.H{"result": randomip.GetPublicIPv6()})
	default:
		message := fmt.Errorf("Invalid version: %s", version)
		c.Error(message)
		c.JSON(http.StatusBadRequest, gin.H{"error": message.Error()})
	}

}
