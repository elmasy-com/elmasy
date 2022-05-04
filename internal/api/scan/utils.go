package scan

import (
	"errors"
	"fmt"
	"sync"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
)

var (
	ErrPortClosed = errors.New("Port is closed")
)

func scanTarget(c chan<- Target, wg *sync.WaitGroup, network, ip, port string) {

	defer wg.Done()

	t := Target{}
	t.Target = ip

	if network == "tcp" {
		open, err := sdk.IsPortOpen("stealth", ip, port)
		if err != nil {
			t.Error = fmt.Errorf("failed to scan %s:%s: %s", ip, port, err)
			c <- t
			return
		}

		if !open {
			t.Error = ErrPortClosed
			c <- t
			return
		}
	}

	tlsResults := make(chan TLS, len(SupportedTLS))

	var twg sync.WaitGroup

	for i := range SupportedTLS {
		twg.Add(1)
		go scanTLS(tlsResults, &twg, SupportedTLS[i], network, ip, port)
	}

	twg.Wait()
	close(tlsResults)

	for i := range tlsResults {
		if t.Error != nil {
			t.Error = i.Error
			c <- t
			return
		}
		t.TLS = append(t.TLS, i)
	}

	c <- t
}

func scanTLS(t chan<- TLS, twg *sync.WaitGroup, version, network, ip, port string) {
	defer twg.Done()

	tls, err := sdk.AnalyzeTLS(version, network, ip, port)
	if err != nil {
		t <- TLS{Error: err}
	}

	t <- TLS{Version: version, Supported: tls.Supported, Ciphers: tls.Ciphers}
}
