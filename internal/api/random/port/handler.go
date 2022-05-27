package port

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(65534) + 1

	switch c.GetHeader("Accept") {
	case "*/*":
		c.String(http.StatusOK, "%d", n)
	case "application/json":
		c.JSON(http.StatusOK, sdk.ResultStr{Result: strconv.Itoa(n)})
	}
}
