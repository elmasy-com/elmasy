package tls

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ssl30"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls10"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls11"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls12"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls13"
	"github.com/elmasy-com/identify"
	"github.com/elmasy-com/slices"
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
	if !identify.IsValidIP(ip) {
		err := fmt.Errorf("Invalid IP: %s", ip)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	port := c.Query("port")
	if !identify.IsValidPort(port) {
		err := fmt.Errorf("Invalid port: %s", port)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	version := c.Query("version")

	switch version {
	case "ssl30":
		r, serr := ssl30.Scan(network, ip+":"+port, 2*time.Second)
		if serr != nil {
			err := fmt.Errorf("Failed to scan ssl30: %s", serr)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"version:": "ssl30", "supported": r.Supported, "ciphers": slices.Strings(r.Ciphers)})

	case "tls10":
		r, serr := tls10.Scan(network, ip+":"+port, 2*time.Second)
		if serr != nil {
			err := fmt.Errorf("Failed to scan tls10: %s", serr)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"version:": "tls10", "supported": r.Supported, "ciphers": slices.Strings(r.Ciphers)})

	case "tls11":
		r, serr := tls11.Scan(network, ip+":"+port, 2*time.Second)
		if serr != nil {
			err := fmt.Errorf("Failed to scan tls11: %s", serr)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"version:": "tls11", "supported": r.Supported, "ciphers": slices.Strings(r.Ciphers)})

	case "tls12":
		r, serr := tls12.Scan(network, ip+":"+port, 2*time.Second)
		if serr != nil {
			err := fmt.Errorf("Failed to scan tls12: %s", serr)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"version:": "tls12", "supported": r.Supported, "ciphers": slices.Strings(r.Ciphers)})

	case "tls13":
		r, serr := tls13.Scan(network, ip, port, 2*time.Second)
		if serr != nil {
			err := fmt.Errorf("Failed to scan tls13: %s", serr)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"version:": "tls13", "supported": r.Supported, "ciphers": slices.Strings(r.Ciphers)})

	default:
		err := fmt.Errorf("Invalid version: %s", version)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
