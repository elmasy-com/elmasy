package main

/*
	Manual testing
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/protocols/tls/tls11"
)

func main() {

	ip := os.Args[1]
	port := os.Args[2]

	r, err := tls11.Scan("tcp", ip+":"+port, 2*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r.Ciphers)
	}
}