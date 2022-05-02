package main

/*
	Manual testing: go run . <ip> <port>
*/
import (
	"fmt"
	"os"
	"time"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/tls12"
)

func main() {

	ip := os.Args[1]
	port := os.Args[2]

	r, err := tls12.Probe("tcp", ip, port, 2*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}
