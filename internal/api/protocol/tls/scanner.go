package tls

import (
	"fmt"
	"strings"
	"time"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ssl30"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls10"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls11"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls12"
	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls13"
	"github.com/elmasy-com/slices"
)

// scanSingle scan TLS using pkg/protocols package.
func scanSingle(version, network, ip, port, servername string) (Result, error) {

	var (
		res Result
		err error
	)

	ip = utils.IPv6BracketAdd(ip)

	switch version {
	case "ssl30":
		var r ssl30.SSL30
		r, err = ssl30.Scan(network, ip, port, 2*time.Second)

		res = Result{IP: ip, Version: "ssl30", Supported: r.Supported, Ciphers: slices.Strings(r.Ciphers)}

	case "tls10":
		var r tls10.TLS10
		r, err = tls10.Scan(network, ip, port, 2*time.Second, tls10.Opts{ServerName: servername})

		res = Result{IP: ip, Version: "tls10", Supported: r.Supported, Ciphers: slices.Strings(r.Ciphers)}

	case "tls11":
		var r tls11.TLS11
		r, err = tls11.Scan(network, ip, port, 2*time.Second, tls11.Opts{ServerName: servername})

		res = Result{IP: ip, Version: "tls11", Supported: r.Supported, Ciphers: slices.Strings(r.Ciphers)}

	case "tls12":
		var r tls12.TLS12
		r, err = tls12.Scan(network, ip, port, 2*time.Second, tls12.Opts{ServerName: servername})

		res = Result{IP: ip, Version: "tls12", Supported: r.Supported, Ciphers: slices.Strings(r.Ciphers)}

	case "tls13":
		var r tls13.TLS13
		r, err = tls13.Scan(network, ip, port, 2*time.Second, tls13.Opts{ServerName: servername})

		res = Result{IP: ip, Version: "tls13", Supported: r.Supported, Ciphers: slices.Strings(r.Ciphers)}

	default:
		return Result{}, fmt.Errorf("Invalid version: %s", version)
	}

	return res, err
}

// scanMany accept target as a domain name, do a DNS lookup and iterate over the ips using the SDK (callback to scanSingle)
func scanMany(version, network, domain, port, servername string) ([]Result, error) {

	targets, err := utils.Lookup46(domain)
	if err != nil {
		return nil, err
	}

	res := make([]Result, 0)

	for i := range targets {

		r, err := sdk.AnalyzeTLS(version, network, strings.Trim(targets[i], "[]"), port, servername)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan %s:%s: %s", targets[i], port, err)
		}

		if len(r) != 1 {
			return nil, fmt.Errorf("Failed to scan: unknwon number of result: %d", len(r))
		}

		res = append(res, Result{IP: targets[i], Version: version, Supported: r[0].Supported, Ciphers: r[0].Ciphers})
	}

	return res, nil
}
