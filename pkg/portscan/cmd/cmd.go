package main

import (
	"fmt"
	"time"

	"github.com/elmasy-com/elmasy/pkg/portscan"
)

/*
Stealth scan requires root.
This small script is for manual testing.
*/

func main() {

	ip := "45.33.32.156" //scanme.nmap.org

	ports := make([]int, 0)
	for i := 1; i < 1025; i++ {
		ports = append(ports, i)
	}

	fmt.Printf("-------------------- stealth --------------------\n")

	t := time.Now()

	r, err := portscan.StealthScan(ip, ports, 1*time.Second)
	if err != nil {
		fmt.Printf("errors: %s\n", err)
	}

	fmt.Printf("Time to scan: %s\n", time.Now().Sub(t).String())

	//fmt.Printf("Result: %#v\n", r)

	fmt.Printf("Open: %d\n", r.Len(portscan.OPEN))
	fmt.Printf("Closed: %#v\n", r.Len(portscan.CLOSED))
	fmt.Printf("Filtered: %#v\n", r.Len(portscan.FILTERED))

	fmt.Printf("-------------------- connect --------------------\n")

	t = time.Now()

	c, err := portscan.ConnectScan(ip, ports, 1*time.Second)
	if err != nil {
		fmt.Printf("errors: %s\n", err)
	}

	fmt.Printf("Time to scan: %s\n", time.Now().Sub(t).String())

	fmt.Printf("Open: %d\n", c.Len(portscan.OPEN))
	fmt.Printf("Closed: %d\n", c.Len(portscan.CLOSED))
	fmt.Printf("Filtered: %d\n", c.Len(portscan.FILTERED))

	fmt.Printf("---------------------- UDP ----------------------\n")

	t = time.Now()

	u, err := portscan.UDPScan(ip, ports, 1*time.Second)
	if err != nil {
		fmt.Printf("errors: %s\n", err)
	}

	fmt.Printf("Time to scan: %s\n", time.Now().Sub(t).String())

	fmt.Printf("Open: %d\n", u.Len(portscan.OPEN))
	fmt.Printf("Closed: %d\n", u.Len(portscan.CLOSED))
	fmt.Printf("Filtered: %d\n", u.Len(portscan.FILTERED))

}
