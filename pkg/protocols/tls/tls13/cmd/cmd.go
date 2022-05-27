package main

/*
	Manual testing: go run . <ip> <port>
*/
import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls13"
)

func main() {

	ip := "[2a01:4f9:c010:81b5::1]"
	port := "443"

	r, err := tls13.Probe("tcp", ip, port, 2*time.Second, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail without SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}

	r, err = tls13.Probe("tcp", ip, port, 2*time.Second, "danielgorbe.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail with invalid SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}

	r, err = tls13.Probe("tcp", ip, port, 2*time.Second, "elmasy.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail with valid SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}
