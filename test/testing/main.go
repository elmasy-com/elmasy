package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/elmasy-com/elmasy/pkg/go-sdk"
)

var OPEN, CLOSED int

func IsPortOpenStealth(ip, port string) (bool, error) {

	ports, errs := sdk.PortScan("stealth", ip, port)
	if errs != nil {

		e := ""

		for i := range errs {
			e += errs[i].Error()
			e += " | "
		}
		return false, fmt.Errorf(e)
	}

	if len(ports) != 1 {
		return false, fmt.Errorf("invalid number of ports: %d", len(ports))
	}

	return ports[0].State == "open", nil

}

func IsPortOpenConnect(ip, port string) (bool, error) {

	ports, errs := sdk.PortScan("connect", ip, port)
	if errs != nil {

		e := ""

		for i := range errs {
			e += errs[i].Error()
			e += " | "
		}
		return false, fmt.Errorf(e)
	}

	if len(ports) != 1 {
		return false, fmt.Errorf("invalid number of ports: %d", len(ports))
	}

	return ports[0].State == "open", nil

}

func IsPortOpen(ip, port string) (bool, error) {

	stealth, err := IsPortOpenStealth(ip, port)
	if err != nil {
		return false, fmt.Errorf("stealth fail: %s", err)
	}

	connect, err := IsPortOpenConnect(ip, port)
	if err != nil {
		return false, fmt.Errorf("connect fail: %s", err)
	}

	if stealth != connect {
		panic(fmt.Sprintf("STEALTH AND CONNECT RESULT ARE DIFFERS: %v vs %v", stealth, connect))
	}

	return stealth, nil
}

// Returns SSL30, TLS10, TLS11, TLS12 support and error.
func ProbeTLS(ip, port string) (bool, bool, bool, bool, error) {

	ssl30, err := sdk.Probe("ssl30", "tcp", ip, port)
	if err != nil {
		return false, false, false, false, fmt.Errorf("failed to check SSL30: %s", err)
	}

	tls10, err := sdk.Probe("tls10", "tcp", ip, port)
	if err != nil {
		return false, false, false, false, fmt.Errorf("failed to check TLS10: %s", err)
	}

	tls11, err := sdk.Probe("tls11", "tcp", ip, port)
	if err != nil {
		return false, false, false, false, fmt.Errorf("failed to check TLS11: %s", err)
	}

	tls12, err := sdk.Probe("tls12", "tcp", ip, port)
	if err != nil {
		return false, false, false, false, fmt.Errorf("failed to check TLS12: %s", err)
	}

	return ssl30, tls10, tls11, tls12, nil
}

func Check(wg *sync.WaitGroup) {

	sdk.DefaultClient.CloseIdleConnections()

	defer wg.Done()

	randomip, err := sdk.GetRandomIP("4")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get randomip: %s\n", err)
		return
	}

	isopen, err := IsPortOpen(randomip, "443")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to portscan: %s\n", err)
		return
	}

	if isopen {
		OPEN++
		ssl30, tls10, tls11, tls12, err := ProbeTLS(randomip, "443")
		if err != nil {
			fmt.Printf("%s:%s is OPEN, but failed to check TLS: %s\n", randomip, "443", err)
		} else {
			fmt.Printf("%s:%s is OPEN, TLS support: ssl30: %v | tls10: %v | tls11: %v | tls12: %v\n", randomip, "443", ssl30, tls10, tls11, tls12)
		}
	} else {
		CLOSED++
		fmt.Printf("%s:%s is closed\n", randomip, "443")
	}
}

func main() {

	var wg sync.WaitGroup

	for a := 0; a < 10; a++ {

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go Check(&wg)
		}

		wg.Wait()

		fmt.Printf("-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------\n")

		// tlsurl := "https://elmasy.com/api/protocol/tls?version=%s&network=tcp&ip=%s&port=%s"

		// ssl30 := TLSScan{}

		// FastGet(fmt.Sprintf(tlsurl, "ssl30", randomip.IP, "443"), &ssl30)

		// fmt.Printf("%#v\n", ssl30)
	}

	fmt.Printf("OPEN PORTS:   %d\n", OPEN)
	fmt.Printf("closed ports: %d\n", CLOSED)

}

/*

Failed to portscan: stealth fail: unknown TCP packet in evaluate(): layers.TCP{BaseLayer:layers.BaseLayer{Contents:[]uint8{0x1, 0xbb, 0xf1, 0x8e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x50, 0x4, 0x0, 0x0, 0x7a, 0x8b, 0x0, 0x0}, Payload:[]uint8{}}, SrcPort:0x1bb, DstPort:0xf18e, Seq:0x0, Ack:0x0, DataOffset:0x5, FIN:false, SYN:false, RST:true, PSH:false, ACK:false, URG:false, ECE:false, CWR:false, NS:false, Window:0x0, Checksum:0x7a8b, Urgent:0x0, sPort:[]uint8{0x1, 0xbb}, dPort:[]uint8{0xf1, 0x8e}, Options:[]layers.TCPOption{}, Padding:[]uint8(nil), opts:[4]layers.TCPOption{layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}, layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}, layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}, layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}}, tcpipchecksum:layers.tcpipchecksum{pseudoheader:layers.tcpipPseudoHeader(nil)}} |

Failed to portscan: stealth fail: unknown TCP packet in evaluate(): layers.TCP{BaseLayer:layers.BaseLayer{Contents:[]uint8{0x1, 0xbb, 0xc0, 0x44, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x50, 0x4, 0x0, 0x0, 0x75, 0x5e, 0x0, 0x0}, Payload:[]uint8{}}, SrcPort:0x1bb, DstPort:0xc044, Seq:0x0, Ack:0x0, DataOffset:0x5, FIN:false, SYN:false, RST:true, PSH:false, ACK:false, URG:false, ECE:false, CWR:false, NS:false, Window:0x0, Checksum:0x755e, Urgent:0x0, sPort:[]uint8{0x1, 0xbb}, dPort:[]uint8{0xc0, 0x44}, Options:[]layers.TCPOption{}, Padding:[]uint8(nil), opts:[4]layers.TCPOption{layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}, layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}, layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}, layers.TCPOption{OptionType:0x0, OptionLength:0x0, OptionData:[]uint8(nil)}}, tcpipchecksum:layers.tcpipchecksum{pseudoheader:layers.tcpipPseudoHeader(nil)}} |

*/
