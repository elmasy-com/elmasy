package portscan

/*
Read more here: https://nmap.org/book/synscan.html
*/
import (
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type stealthOpts struct {
	network   string                // Must be "ip4:tcp" or "ip6:tcp"
	timeout   time.Duration         // Read timeout
	result    Result                //
	raddr     *net.IPAddr           // Remote address
	dconn     *net.IPConn           // Connection to send packets
	laddr     *net.IPAddr           // Local address
	lport     uint16                // Source port of outgoing packets
	lconn     *net.IPConn           // COnnection to read pakcets
	nlayer    gopacket.NetworkLayer // Network layer used to create tcp packet
	fldecoder gopacket.Decoder      // First layer decoder used in read
	m         sync.Mutex
}

// newStealthOpts returns a StealthOpts struct.
func newStealthOpts(ip string, ports []int, timeout time.Duration) (*stealthOpts, error) {

	var (
		scan stealthOpts
		err  error
	)

	// Check valid IPv4/IPv6 address and set network
	switch i := net.ParseIP(ip); true {
	case i == nil:
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	case i.To4() != nil:
		scan.network = "ip4:tcp"
	case i.To16() != nil:
		scan.network = "ip6:tcp"
	default:
		return nil, fmt.Errorf("unparsable IP address: %s", ip)

	}

	scan.timeout = timeout

	// Set result
	for i := range ports {

		if ports[i] < 1 || ports[i] > 65535 {
			return nil, fmt.Errorf("invalid port: %d\n", ports[i])
		}

		scan.result = append(scan.result, Port{Port: ports[i], State: FILTERED})
	}

	// Set raddr
	if scan.raddr, err = net.ResolveIPAddr(scan.network, ip); err != nil {
		return nil, fmt.Errorf("failed to resolve remote address: %s", err)
	}

	// Set dconn
	scan.dconn, err = net.DialIP(scan.network, nil, scan.raddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	// Set laddr
	if scan.laddr, err = net.ResolveIPAddr(scan.network, scan.dconn.LocalAddr().String()); err != nil {
		scan.dconn.Close()
		return nil, fmt.Errorf("failed to resolve local address: %s", err)
	}

	// Set lport
	rand.Seed(time.Now().UnixNano())
	scan.lport = uint16(rand.Intn(16383) + 49152)

	// Set lconn
	scan.lconn, err = net.ListenIP(scan.network, scan.laddr)
	if err != nil {
		scan.dconn.Close()
		return nil, fmt.Errorf("failed to listen: %s", err)
	}

	// Set read initial read timeout.
	// The timeout will be updated after every packet from the target.
	if err = scan.lconn.SetDeadline(time.Now().Add(scan.timeout)); err != nil {
		scan.dconn.Close()
		scan.lconn.Close()
		return nil, fmt.Errorf("failed to set initial deadline: %s", err)
	}

	// Set nlayer
	switch scan.network {
	case "ip4:tcp":
		scan.nlayer = &layers.IPv4{SrcIP: scan.laddr.IP, DstIP: scan.raddr.IP}
	case "ip6:tcp":
		scan.nlayer = &layers.IPv6{SrcIP: scan.laddr.IP, DstIP: scan.raddr.IP}
	}

	// Set fldecoder
	// IPv4 raw socket read() returns the IP header, but Go's ReadFromIP() returns only the TCP header
	scan.fldecoder = layers.LayerTypeTCP

	return &scan, nil
}

func (s *stealthOpts) close() {
	s.dconn.Close()
	s.lconn.Close()
}

// create the tcp packets and send into c to send()
func (s *stealthOpts) createPackets(c chan<- []byte, e chan<- error, wg *sync.WaitGroup) {

	defer wg.Done()
	defer close(c)

	for i := range s.result {

		tcp := layers.TCP{
			SrcPort: layers.TCPPort(s.lport),
			DstPort: layers.TCPPort(s.result[i].Port),
			Seq:     rand.Uint32(),
			Window:  64240,
			SYN:     true,
			Options: []layers.TCPOption{
				{OptionType: 2, OptionLength: 4, OptionData: []byte{0x05, 0xb4}}, // MSS=1460
			},
		}

		tcp.SetNetworkLayerForChecksum(s.nlayer)

		buf := gopacket.NewSerializeBuffer()

		opts := gopacket.SerializeOptions{
			ComputeChecksums: true,
			FixLengths:       true,
		}

		if err := gopacket.SerializeLayers(buf, opts, &tcp); err != nil {
			e <- fmt.Errorf("failed to serialize for port %d: %s", s.result[i].Port, err)
			break
		}

		c <- buf.Bytes()
	}
}

// Send TCP packets get from createPackets through c channel
func (s *stealthOpts) send(c <-chan []byte, e chan<- error, wg *sync.WaitGroup) {

	defer wg.Done()

	for p := range c {

		if _, err := s.dconn.Write(p); err != nil {
			e <- fmt.Errorf("failed to write: %s", err)
		}
	}
}

// Try to read the response, assemble the packet and send it into c.
func (s *stealthOpts) read(c chan<- layers.TCP, e chan<- error, wg *sync.WaitGroup) {

	defer wg.Done()
	defer close(c)

	for {

		// Stop reading if every s.result if filled
		s.m.Lock()
		if s.result.Len(FILTERED) == 0 {
			break
		}
		s.m.Unlock()

		b := make([]byte, 256)

		n, raddr, err := s.lconn.ReadFromIP(b)
		if err != nil {
			if !strings.Contains(err.Error(), "i/o timeout") {
				e <- fmt.Errorf("fail in read(): %s\n", err)
			}
			break
		}

		// Skip packets that not comes from the target IP.
		if !s.raddr.IP.Equal(raddr.IP) {
			continue
		}

		b = b[:n]

		packet := gopacket.NewPacket(b, s.fldecoder, gopacket.Default)

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			e <- fmt.Errorf("fail in read(): failed to convert TCP layer\n")
			continue
		}

		c <- *tcp
	}
}

// Read TCP packets from channel c, and evalute the state based on the TCP flags.
func (s *stealthOpts) evaluate(c <-chan layers.TCP, e chan<- error, wg *sync.WaitGroup) {

	defer wg.Done()

	for t := range c {

		if uint16(t.DstPort) != s.lport {
			continue
		}

		s.m.Lock()

		switch {
		case t.SYN && t.ACK && !t.RST && !t.FIN && !t.PSH && !t.URG:
			// Normal SYN/ACK, the port is open
			s.result.addResult(int(t.SrcPort), OPEN)

		case !t.SYN && t.ACK && t.RST && !t.FIN && !t.PSH && !t.URG:
			// Normal RST/ACK, the port is closed
			s.result.addResult(int(t.SrcPort), CLOSED)

		case !t.SYN && t.ACK && !t.RST && !t.FIN && !t.PSH && !t.URG && t.Window == 0:
			// TCP ZeroWindow (see: https://osqa-ask.wireshark.org/questions/2365/tcp-window-size-and-scaling/)
			s.result.addResult(int(t.SrcPort), CLOSED)

		default:
			e <- fmt.Errorf("unknown TCP packet in evaluate(): %#v", t)
		}

		// Update the read deadline
		s.lconn.SetDeadline(time.Now().Add(s.timeout))

		s.m.Unlock()
	}
}

func synScan(scan *stealthOpts) (Result, []error) {

	packets := make(chan []byte, len(scan.result))
	tcps := make(chan layers.TCP, len(scan.result))
	errc := make(chan error, len(scan.result))
	var errs []error
	var wg sync.WaitGroup

	wg.Add(4)

	go scan.read(tcps, errc, &wg)
	go scan.evaluate(tcps, errc, &wg)

	go scan.createPackets(packets, errc, &wg)
	go scan.send(packets, errc, &wg)

	wg.Wait()
	close(errc)

	for e := range errc {
		errs = append(errs, e)
	}

	return scan.result, errs
}

// StealthScan uses syn scan method to scan ports.
// ip must be a valid IPv4/IPv6 address, not a domain.
// StealthScan retries the FILTERED ports for once again.
func StealthScan(ip string, ports []int, timeout time.Duration) (Result, []error) {

	result := make(Result, 0)

	scan1, err := newStealthOpts(ip, ports, timeout)
	if err != nil {
		return nil, []error{err}
	}

	scan1r, errs := synScan(scan1)
	scan1.close()
	if errs != nil {
		return nil, errs
	}

	result = append(result, scan1r.GetPorts(OPEN)...)
	result = append(result, scan1r.GetPorts(CLOSED)...)

	// Scan filtered ports again, maybe a network error
	if scan1r.Len(FILTERED) > 0 {
		scan2, err := newStealthOpts(ip, scan1r.GetPortsInt(FILTERED), timeout)
		if err != nil {
			return nil, []error{err}
		}

		scan2r, errs := synScan(scan2)
		scan2.close()
		if errs != nil {
			return nil, errs
		}

		result = append(result, scan2r...)
	}

	// Sort result by port number
	sort.Slice(result, func(i, j int) bool { return result[i].Port < result[j].Port })

	return result, nil
}
