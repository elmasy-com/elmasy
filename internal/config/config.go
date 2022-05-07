package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"gopkg.in/yaml.v2"
)

type TLS struct {
	Cert string
	Key  string
}

type Config struct {
	Verbose             bool     `yaml:"verbose"`
	URL                 string   `yaml:"url"`
	Listen              string   `yaml:"listen"`
	TrustedProxies      []string `yaml:"trusted_proxies"`
	BlackListedNetsStr  []string `yaml:"blacklisted_networks"`
	SSLCertificate      string   `yaml:"ssl_certificate"`
	SSLKey              string   `yaml:"ssl_certificate_key"`
	BlacklistedNetworks []net.IPNet
}

var GlobalConfig Config

func ParseConfig(path string) error {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %s", path, err)
	}

	if err := yaml.Unmarshal(bytes, &GlobalConfig); err != nil {
		return err
	}

	if err := GlobalConfig.convertBlackListedIPs(); err != nil {
		return err
	}

	// Default URL
	if GlobalConfig.URL == "" {
		GlobalConfig.URL = "https://elmasy.com"
	} else {
		GlobalConfig.URL = strings.TrimRight(GlobalConfig.URL, "/")
	}
	sdk.API_PATH = GlobalConfig.URL + "/api"

	if GlobalConfig.Verbose {
		fmt.Printf("Config:\n%#v\n\n", GlobalConfig)
	}

	return nil
}

func (c *Config) convertBlackListedIPs() error {

	for i := range c.BlackListedNetsStr {

		_, net, err := net.ParseCIDR(c.BlackListedNetsStr[i])
		if err != nil {
			return err
		}
		c.BlacklistedNetworks = append(c.BlacklistedNetworks, *net)
	}

	return nil
}
