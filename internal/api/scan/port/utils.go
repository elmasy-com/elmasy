package port

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/elmasy-com/elmasy/internal/config"
	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

type Params struct {
	Technique string
	IP        string
	Timeout   time.Duration
	Port      int
}

type PortString struct {
	Port  string `json:"port"`
	State string `json:"state"`
}

// Parse the query params and handle error.
// In case of error, this function set error in the context and also return it.
// So, in the calling function, it is enough to return if err != nil.
func parseQuery(c *gin.Context) (Params, error) {

	params := Params{}

	params.Technique = c.DefaultQuery("technique", "connect")

	if params.Technique != "syn" &&
		params.Technique != "stealth" &&
		params.Technique != "connect" {

		err := fmt.Errorf("Invalid technique: %s", params.Technique)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}

	params.IP = c.Query("ip")
	if params.IP == "" {
		err := fmt.Errorf("ip is missing")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}
	if !identify.IsValidIP(params.IP) {
		err := fmt.Errorf("Invalid IP address: %s", params.IP)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}
	if isIPBlacklisted(params.IP) {
		err := fmt.Errorf("Blacklisted IP address: %s", params.IP)
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return params, err
	}
	if identify.IsValidIPv6(params.IP) {
		params.IP = utils.IPv6BracketAdd(params.IP)
	}

	timeoutQuery := c.DefaultQuery("timeout", "2")
	timeoutInt, err := strconv.Atoi(timeoutQuery)
	if err != nil || timeoutInt < 0 {
		err := fmt.Errorf("Invalid timeout: %s", timeoutQuery)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}
	params.Timeout = time.Duration(timeoutInt) * time.Second

	portQuery := c.Query("port")
	if portQuery == "" {
		err := fmt.Errorf("port is missing")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}
	params.Port, err = strconv.Atoi(portQuery)
	if err != nil || !identify.IsValidPort(params.Port) {
		err := fmt.Errorf("Invalid port: %s", portQuery)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}

	return params, nil
}

// IsIPBlacklisted decides whether ip is in the blacklisted ip address range.
func isIPBlacklisted(ip string) bool {

	p := net.ParseIP(ip)
	if p == nil {
		panic("Failed to parse IP in isIPBlacklisted(): " + ip)
	}

	for i := range config.GlobalConfig.BlacklistedNetworks {

		if config.GlobalConfig.BlacklistedNetworks[i].Contains(p) {
			return true
		}
	}

	return false
}
