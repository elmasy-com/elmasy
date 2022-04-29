package dns

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

// Probe checks whether DNS protocol is available on network on ip:port.
// This function do a simple query with "elmasy.com."/"A".
// network must be "tcp", "tcp-tls" or "udp".
func Probe(network, ip, port string, timeout time.Duration) (bool, error) {

	if network != "tcp" && network != "tcp-tls" && network != "udp" {
		return false, fmt.Errorf("invalid network: %s", network)
	}

	m := new(dns.Msg)
	m.SetQuestion("elmasy.com.", dns.TypeA)

	c := new(dns.Client)

	c.Net = network
	c.Timeout = timeout

	_, _, err := c.Exchange(m, ip+":"+port)

	return err == nil, nil
}
