package main

import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls10"
)

func main() {

	ip := os.Args[1]
	port := os.Args[2]

	r, err := tls10.Scan("tcp", ip+":"+port, 2*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r.Ciphers)
	}
}
