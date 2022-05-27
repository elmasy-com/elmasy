package main

/*
	Manual testing
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls11"
)

func main() {

	ip := "142.132.164.231"
	port := "443"

	r, err := tls11.Probe("tcp", ip, port, 2*time.Second, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail without SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}

	r, err = tls11.Probe("tcp", ip, port, 2*time.Second, "danielgorbe.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail with invalid SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}

	r, err = tls11.Probe("tcp", ip, port, 2*time.Second, "elmasy.com")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail with valid SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}
