package scan

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

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

	ips := make([]string, 0)
	var servername string

	if identify.IsDomainName(target) {
		ips, err = utils.Lookup46(target)
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

		servername = target
	} else {
		ips = append(ips, target)
	}

	targets := make(chan sdk.Target, len(ips))
	errChan := make(chan error, 20)
	var wg sync.WaitGroup

	for i := range ips {
		wg.Add(1)
		go scanTarget(targets, errChan, &wg, network, ips[i], port, servername)
	}

	wg.Wait()
	close(targets)
	close(errChan)

	result := sdk.Result{Domain: target}

	for e := range errChan {
		c.Error(e)
		result.Errors = append(result.Errors, e.Error())
	}

	for t := range targets {

		result.Targets = append(result.Targets, t)

	}

	c.JSON(http.StatusOK, result)

}
