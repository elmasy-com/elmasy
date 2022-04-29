package probe

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/dns"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ssl30"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls10"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls11"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls12"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	network := c.DefaultQuery("network", "tcp")
	if network != "tcp" && network != "udp" {
		err := fmt.Errorf("Invalid network: %s", network)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ip := c.Query("ip")
	if ip == "" || !identify.IsValidIP(ip) {
		err := fmt.Errorf("Invalid IP: %s", ip)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	port := c.Query("port")
	if port == "" || !identify.IsValidPort(port) {
		err := fmt.Errorf("Invalid port: %s", port)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		supported bool
		err       error
	)
	switch protocol := c.Query("protocol"); protocol {
	case "dns":
		supported, err = dns.Probe(network, ip, port, 2*time.Second)
	case "ssl30":
		supported, err = ssl30.Probe(network, ip, port, 2*time.Second)
	case "tls10":
		supported, err = tls10.Probe(network, ip, port, 2*time.Second)
	case "tls11":
		supported, err = tls11.Probe(network, ip, port, 2*time.Second)
	case "tls12":
		supported, err = tls12.Probe(network, ip, port, 2*time.Second)
	default:
		err = fmt.Errorf("Invalid protocol: %s", protocol)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": supported})

}
