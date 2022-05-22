package tls

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

type Result struct {
	IP        string   `json:"ip"`
	Version   string   `json:"version"`
	Supported bool     `json:"supported"`
	Ciphers   []string `json:"ciphers"`
}

func Get(c *gin.Context) {

	version := c.Query("version")
	if version == "" {
		err := fmt.Errorf("version is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if version != "ssl30" && version != "tls10" && version != "tls11" && version != "tls12" && version != "tls13" {
		err := fmt.Errorf("Invalid version: %s", version)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	servername := c.Query("servername")

	var (
		result = make([]Result, 0)
		err    error
	)

	switch true {
	case identify.IsValidIP(target):
		var r Result
		r, err = scanSingle(version, network, target, port, getServerName(servername, target))

		result = append(result, r)
	case identify.IsDomainName(target):
		result, err = scanMany(version, network, target, port, getServerName(servername, target))
	default:
		err = fmt.Errorf("invalid target: %s", target)
	}

	if err != nil {
		err := fmt.Errorf("Failed to scan: %s", err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
