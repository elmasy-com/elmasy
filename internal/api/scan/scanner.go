package scan

import (
	"fmt"
	"sort"
	"sync"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
)

var SupportedTLS = []string{"ssl30", "tls10", "tls11", "tls12", "tls13"}

func scanTarget(c chan<- sdk.Target, errors chan<- error, wg *sync.WaitGroup, network, ip, port, servername string) {

	defer wg.Done()

	t := sdk.Target{IP: ip}

	if network == "tcp" {
		state, err := sdk.PortScan("connect", ip, port, "2")
		if err != nil {
			errors <- fmt.Errorf("failed to scan %s:%s: %s", ip, port, err)
			c <- t
			return
		}

		if state != "open" {
			errors <- fmt.Errorf("%s:%s is %s", ip, port, state)
			c <- t
			return
		}
	}

	tlsResults := make(chan sdk.TLSVersion, len(SupportedTLS))
	certResult := make(chan sdk.Cert, 1)
	var twg sync.WaitGroup

	for i := range SupportedTLS {
		twg.Add(1)
		go scanTLS(tlsResults, errors, &twg, SupportedTLS[i], network, ip, port, servername)
	}

	twg.Add(1)
	scanTLSCert(certResult, errors, &twg, network, ip, port, servername)

	twg.Wait()
	close(tlsResults)
	close(certResult)

	for i := range tlsResults {

		t.TLS.Versions = append(t.TLS.Versions, i)
	}

	t.TLS.Cert = <-certResult

	sort.Slice(t.TLS.Versions, func(i, j int) bool { return t.TLS.Versions[i].Version < t.TLS.Versions[j].Version })

	c <- t
}

func scanTLS(t chan<- sdk.TLSVersion, errors chan<- error, twg *sync.WaitGroup, version, network, ip, port, servername string) {
	defer twg.Done()

	tls, err := sdk.AnalyzeTLS(version, network, ip, port, servername)

	if err != nil {
		errors <- fmt.Errorf("%s://%s:%s -> TLS %s: %s", network, utils.IPv6BracketAdd(ip), port, version, err)
		return
	}

	t <- sdk.TLSVersion{Version: version, Supported: tls.Supported, Ciphers: tls.Ciphers}
}

func scanTLSCert(t chan<- sdk.Cert, errors chan<- error, twg *sync.WaitGroup, network, ip, port, servername string) {
	defer twg.Done()

	cert, err := sdk.GetCertificate(network, ip, port, servername)
	if err != nil {
		errors <- fmt.Errorf("%s://%s:%s -> TLS Cert: %s", network, utils.IPv6BracketAdd(ip), port, err)
		return
	}

	if cert.VerifiedError != "" {
		errors <- fmt.Errorf("%s://%s:%s -> TLS Cert: %s", network, utils.IPv6BracketAdd(ip), port, cert.VerifiedError)
	}

	t <- cert
}
