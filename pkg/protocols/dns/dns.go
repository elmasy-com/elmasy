package dns

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/elmasy-com/identify"
	"github.com/miekg/dns"
)

var conf *dns.ClientConfig

// The number of tries in case of error.
var MAX_RETRIES int = 2

func init() {

	var err error

	conf, err = dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		panic("Failed to parse /etc/resolv.conf: " + err.Error())
	}

	rand.Seed(time.Now().UnixMicro())
}

// getServer used to randomize DNS servers.
func getServer() string {

	var r string

	if len(conf.Servers) == 1 {
		r = conf.Servers[0]
	} else {
		r = conf.Servers[rand.Intn(len(conf.Servers))]
	}

	if identify.IsValidIPv6(r) {
		r = "[" + r + "]"
	}

	return r + ":53"
}

// Generic query for internal use.
// Returns the Answer section.
// In case of error, returns nil.
func query(name string, t uint16) ([]dns.RR, error) {

	name = dns.Fqdn(name)

	msg := new(dns.Msg)
	msg.SetQuestion(name, t)

	in, err := dns.Exchange(msg, getServer())

	if err != nil {
		return nil, err
	}

	if in.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf(dns.RcodeToString[in.Rcode])
	}

	return in.Answer, nil
}
