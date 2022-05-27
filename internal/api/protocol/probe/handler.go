package probe

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/elmasy/pkg/protocols/dns"
	etls "github.com/elmasy-com/elmasy/pkg/protocols/tls"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	network := c.DefaultQuery("network", "tcp")
	if network != "tcp" && network != "udp" {
		err := fmt.Errorf("Invalid network: %s", network)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	ip := c.Query("ip")
	if ip == "" {
		err := fmt.Errorf("ip is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}
	if !identify.IsValidIP(ip) {
		err := fmt.Errorf("Invalid ip: %s", ip)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	port := c.Query("port")
	if port == "" {
		err := fmt.Errorf("port is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}
	if !identify.IsValidPort(port) {
		err := fmt.Errorf("Invalid port: %s", port)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	var (
		supported bool
		err       error
	)

	switch protocol := c.Query("protocol"); protocol {
	case "dns":
		supported, err = dns.Probe(network, utils.IPv6BracketAdd(ip), port, 2*time.Second)
	case "ssl30", "tls10", "tls11", "tls12", "tls13":
		supported, err = etls.Probe(protocol, network, utils.IPv6BracketAdd(ip), port, 2*time.Second, "")
	default:
		err = fmt.Errorf("Invalid protocol: %s", protocol)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, sdk.Error{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sdk.ResultBool{Result: supported})

}
