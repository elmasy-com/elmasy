package port

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/elmasy-com/elmasy/internal/config"
	"github.com/elmasy-com/elmasy/pkg/portscan"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

type Params struct {
	Technique string
	IP        string
	Ports     []int
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

	portsQuery := c.Query("ports")
	if portsQuery == "" {
		err := fmt.Errorf("ports is missing")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return params, err
	}

	portsSlice := strings.Split(portsQuery, ",")
	for i := range portsSlice {

		p, err := strconv.Atoi(portsSlice[i])
		if err != nil {
			err := fmt.Errorf("Invalid port: \"%s\"", portsSlice[i])
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return params, err
		}
		if !identify.IsValidPort(p) {
			err := fmt.Errorf("Invalid port: %d", p)
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return params, err
		}
		params.Ports = append(params.Ports, p)
	}

	if len(params.Ports) > 100 {
		err := fmt.Errorf("Too much port: %d", len(params.Ports))
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

	for i := range config.GlobalConfig.BlacklistedNets {

		if config.GlobalConfig.BlacklistedNets[i].Contains(p) {
			return true
		}
	}

	return false
}

// Convert portscan.Port's int/int to string/string for better result of the API.
func convertResultString(r portscan.Result) []PortString {

	ps := make([]PortString, 0)

	for i := range r {
		ps = append(ps, PortString{Port: strconv.Itoa(r[i].Port), State: r[i].State.String()})
	}

	return ps
}
