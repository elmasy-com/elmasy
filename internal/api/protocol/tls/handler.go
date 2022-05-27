package tls

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	etls "github.com/elmasy-com/elmasy/pkg/protocols/tls"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

type Cipher struct {
	Name     string
	Security string
}

type Result struct {
	IP        string   `json:"ip"`
	Version   string   `json:"version"`
	Supported bool     `json:"supported"`
	Ciphers   []Cipher `json:"ciphers"`
}

func Get(c *gin.Context) {

	version := c.Query("version")
	if version == "" {
		err := fmt.Errorf("version is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}
	if version != "ssl30" && version != "tls10" && version != "tls11" && version != "tls12" && version != "tls13" {
		err := fmt.Errorf("Invalid version: %s", version)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

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
	ip = utils.IPv6BracketAdd(ip)

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

	servername := c.Query("servername")

	r, err := etls.Scan(version, network, ip, port, 2*time.Second, servername)
	if err != nil {
		err := fmt.Errorf("Failed to scan: %s", err)
		c.Error(err)
		c.JSON(http.StatusInternalServerError, sdk.Error{Err: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sdk.TLSVersion{Version: version, Supported: r.Supported, Ciphers: convertCiphers(r.Ciphers)})
}
