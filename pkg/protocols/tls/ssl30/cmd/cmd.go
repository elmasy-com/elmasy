package main

/*
	For manual testing.
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ssl30"
)

func main() {

	ip := os.Args[1]
	port := os.Args[2]

	r, err := ssl30.Probe("tcp", ip, port, 2*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}
