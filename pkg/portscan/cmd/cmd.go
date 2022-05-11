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
		fmt.Printf("error: %s\n", err)
	}

	fmt.Printf("Time to scan: %s\n", time.Now().Sub(t).String())

	//fmt.Printf("Result: %#v\n", r)

	fmt.Printf("Open: %#v\n", r.GetPorts(portscan.OPEN))
	//fmt.Printf("Closed: %#v\n", portscan.GetPorts(r, portscan.CLOSED))
	fmt.Printf("Filtered: %#v\n", r.GetPorts(portscan.FILTERED))

	fmt.Printf("-------------------- connect --------------------\n")

	t = time.Now()

	c, err := portscan.ConnectScan(ip, ports, 1*time.Second)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	fmt.Printf("Time to scan: %s\n", time.Now().Sub(t).String())

	//fmt.Printf("Result: %#v\n", r)

	fmt.Printf("Open: %#v\n", c.GetPorts(portscan.OPEN))
	//fmt.Printf("Closed: %#v\n", portscan.GetPorts(c, portscan.CLOSED))
	fmt.Printf("Filtered: %#v\n", c.GetPorts(portscan.FILTERED))

}
