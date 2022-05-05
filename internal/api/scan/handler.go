package scan

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

type TLS struct {
	Version   string   `json:"version"`
	Supported bool     `json:"supported"`
	Ciphers   []string `json:"ciphers"`
	Error     error    `json:"-"`
}

type Target struct {
	Target string `json:"target"`
	TLS    []TLS  `json:"tls"`
	Error  error  `json:"-"`
}

type Result struct {
	Result []Target `json:"result"`
}

var SupportedTLS = []string{"ssl30", "tls10", "tls11", "tls12"}

func Get(c *gin.Context) {

	var err error

	target := c.Query("target")
	if target == "" {
		err = fmt.Errorf("Target is missing")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !identify.IsDomainName(target) && !identify.IsValidIPv4(target) {
		err = fmt.Errorf("Invalid target: %s", target)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	port := c.DefaultQuery("port", "443")
	if !identify.IsValidPort(port) {
		err = fmt.Errorf("Invalid port: %s", port)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	network := c.DefaultQuery("network", "tcp")
	if network != "tcp" && network != "udp" {
		err = fmt.Errorf("Invalid network: %s", network)
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ips []string

	if identify.IsDomainName(target) {
		ips, err = sdk.DNSLookup("A", target)
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
		ips = append(ips, target)
	}

	targets := make(chan Target, len(ips))
	var wg sync.WaitGroup

	for i := range ips {

		wg.Add(1)

		scanTarget(targets, &wg, network, ips[i], port)
	}

	wg.Wait()
	close(targets)

	result := Result{}

	for t := range targets {
		if t.Error != nil {

			if errors.Is(t.Error, ErrPortClosed) {
				c.Error(t.Error)
				c.JSON(http.StatusForbidden, gin.H{"error": t.Error.Error()})
				return
			}

			err = fmt.Errorf("failed to check %s: %s", t.Target, t.Error)
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result.Result = append(result.Result, t)

	}

	c.JSON(http.StatusOK, result.Result)
}
