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

	ip := "142.132.164.231"
	port := "443"

	r, err := tls12.Probe("tcp", ip, port, 2*time.Second, tls12.Opts{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail without SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}

	r, err = tls12.Probe("tcp", ip, port, 2*time.Second, tls12.Opts{ServerName: "danielgorbe.com"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail with invalid SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}

	r, err = tls12.Probe("tcp", ip, port, 2*time.Second, tls12.Opts{ServerName: "elmasy.com"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail with valid SNI: %s\n", err)
	} else {
		fmt.Printf("%#v\n", r)
	}
}
