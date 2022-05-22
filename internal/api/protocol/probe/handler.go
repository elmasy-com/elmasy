package probe

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/protocols/dns"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ssl30"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls10"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls11"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls12"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls13"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

type Probe struct {
	IP        string // `json:"ip"`
	Supported bool   // `json:"supported"`
}

func Get(c *gin.Context) {

	network := c.DefaultQuery("network", "tcp")
	if network != "tcp" && network != "udp" {
		err := fmt.Errorf("Invalid network: %s", network)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target := c.Query("target")
	if target == "" {
		err := fmt.Errorf("target is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !identify.IsValidIP(target) && !identify.IsDomainName(target) {
		err := fmt.Errorf("Invalid target: %s", target)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	port := c.Query("port")
	if port == "" {
		err := fmt.Errorf("port is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !identify.IsValidPort(port) {
		err := fmt.Errorf("Invalid port: %s", port)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targets := make([]string, 0)

	if identify.IsDomainName(target) {
		var err error

		targets, err = utils.Lookup46(target)
		if err != nil {
			var code int

			if err.Error() == "NXDOMAIN" {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}

			err = fmt.Errorf("Lookup failed: %s", err)
			c.Error(err)
			c.JSON(code, gin.H{"error": err.Error()})
			return
		}

	} else {
		targets = append(targets, target)
	}

	var (
		result    = make([]Probe, 0)
		supported bool
		err       error
	)

	for i := range targets {

		switch protocol := c.Query("protocol"); protocol {
		case "dns":
			supported, err = dns.Probe(network, utils.IPv6BracketAdd(targets[i]), port, 2*time.Second)
		case "ssl30":
			supported, err = ssl30.Probe(network, utils.IPv6BracketAdd(targets[i]), port, 2*time.Second)
		case "tls10":
			supported, err = tls10.Probe(network, utils.IPv6BracketAdd(targets[i]), port, 2*time.Second, tls10.Opts{ServerName: target})
		case "tls11":
			supported, err = tls11.Probe(network, utils.IPv6BracketAdd(targets[i]), port, 2*time.Second, tls11.Opts{ServerName: target})
		case "tls12":
			supported, err = tls12.Probe(network, utils.IPv6BracketAdd(targets[i]), port, 2*time.Second, tls12.Opts{ServerName: target})
		case "tsl13":
			supported, err = tls13.Probe(network, utils.IPv6BracketAdd(targets[i]), port, 2*time.Second, tls13.Opts{ServerName: target})
		default:
			err = fmt.Errorf("Invalid protocol: %s", protocol)
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result = append(result, Probe{IP: targets[i], Supported: supported})
	}

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)

}
