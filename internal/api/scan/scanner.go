package scan

import (
	"errors"
	"fmt"
	"sync"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
)

var SupportedTLS = []string{"ssl30", "tls10", "tls11", "tls12", "tls13"}

var (
	ErrPortClosed = errors.New("Port is closed")
)

func scanTarget(c chan<- Target, wg *sync.WaitGroup, network, ip, port, servername string) {

	defer wg.Done()

	t := Target{}
	t.Target = ip

	if network == "tcp" {
		state, err := sdk.PortScan("connect", ip, port, "2")
		if err != nil {
			t.Error = fmt.Errorf("failed to scan %s:%s: %s", ip, port, err)
			c <- t
			return
		}

		if state != "open" {
			t.Error = ErrPortClosed
			c <- t
			return
		}
	}

	tlsResults := make(chan TLS, len(SupportedTLS))

	var twg sync.WaitGroup

	for i := range SupportedTLS {
		twg.Add(1)
		go scanTLS(tlsResults, &twg, SupportedTLS[i], network, ip, port, servername)
	}

	twg.Wait()
	close(tlsResults)

	for i := range tlsResults {
		if i.Error != nil {
			t.Error = i.Error
			c <- t
			return
		}
		t.TLS = append(t.TLS, i)
	}

	c <- t
}

func scanTLS(t chan<- TLS, twg *sync.WaitGroup, version, network, ip, port, servername string) {
	defer twg.Done()

	tls, err := sdk.AnalyzeTLS(version, network, ip, port, servername)

	if err != nil {
		t <- TLS{Error: err}
		return
	}

	if len(tls) != 1 {
		t <- TLS{Error: fmt.Errorf("multiple result for one scan")}
		return
	}

	t <- TLS{Version: version, Supported: tls[0].Supported, Ciphers: tls[0].Ciphers, Error: err}
}
