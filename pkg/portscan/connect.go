package portscan

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type connectOpts struct {
	address string
	timeout time.Duration
	result  Result
	errs    []error
	lock    *semaphore.Weighted
}

func newConnectOpts(address string, ports []int, timeout time.Duration) (connectOpts, error) {

	var scan connectOpts

	scan.address = address

	scan.timeout = timeout

	for i := range ports {
		scan.result = append(scan.result, Port{Port: ports[i], State: FILTERED})
	}

	ulimit, err := GetMaxFD()
	if err != nil {
		return scan, fmt.Errorf("failed to get max file descriptors: %s", err)
	}
	scan.lock = semaphore.NewWeighted(int64(ulimit))

	return scan, nil

}

func (c *connectOpts) connect(p *Port, wg *sync.WaitGroup, e chan<- error) {

	target := fmt.Sprintf("%s:%d", c.address, p.Port)

	//fmt.Printf("Scanning %s...\n", target)

	conn, err := net.DialTimeout("tcp", target, c.timeout)
	if err != nil {

		switch {
		// Possible DROP on firewall
		case strings.Contains(err.Error(), "i/o timeout"):
			p.State = FILTERED
			wg.Done()

		// RST (can be ICMP: Destination unreachable (Port unreachable))
		case strings.Contains(err.Error(), "connection refused"):
			p.State = CLOSED
			wg.Done()

		// ICMP (type: 3/code: 13): Destination unreachable (Communication administratively filtered)
		case strings.Contains(err.Error(), "no route to host"):
			p.State = FILTERED
			wg.Done()

		// ICMP (type: 3/code: 0): Destination unreachable (Network unreachable)
		case strings.Contains(err.Error(), "network is unreachable"):
			p.State = FILTERED
			wg.Done()

		case strings.Contains(err.Error(), "too many open files"):
			time.Sleep(c.timeout)
			c.connect(p, wg, e)
		default:
			e <- fmt.Errorf("unknown error: %s", err)
			wg.Done()

		}
	} else {
		p.State = OPEN
		conn.Close()
		wg.Done()
	}

}

func connectscan(opts connectOpts) (Result, []error) {

	var errs []error
	ec := make(chan error, len(opts.result))
	var wg sync.WaitGroup

	for i := range opts.result {

		wg.Add(1)
		go opts.connect(&opts.result[i], &wg, ec)
	}

	wg.Wait()
	close(ec)

	for e := range ec {
		errs = append(errs, e)
	}

	return opts.result, errs
}

// ConnectScan uses the basic connect() method to determine if ports are open.
// address can be a hostname, not limited to IP address.
// ConnectScan retries the FILTERED ports for once again.
func ConnectScan(address string, ports []int, timeout time.Duration) (Result, []error) {

	var results Result

	scan1, err := newConnectOpts(address, ports, timeout)
	if err != nil {
		return nil, []error{err}
	}

	scan1r, errs := connectscan(scan1)
	if err != nil {
		return nil, errs
	}

	results = append(results, scan1r.GetPorts(OPEN)...)
	results = append(results, scan1r.GetPorts(CLOSED)...)

	if scan1r.Len(FILTERED) > 0 {
		scan2, err := newConnectOpts(address, scan1r.GetPortsInt(FILTERED), timeout)
		if err != nil {
			return nil, []error{err}
		}

		scan2r, errs := connectscan(scan2)
		if errs != nil {
			return nil, errs
		}

		results = append(results, scan2r...)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Port < results[j].Port })

	return results, nil
}
