package main

/*
	For manual testing.
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/protocols/tls/ssl30"
)

func main() {

	ip := os.Args[1]
	port := os.Args[2]

	r, err := ssl30.Scan("tcp", ip+":"+port, 2*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}
