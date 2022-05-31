package portscan

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type udpOpts struct {
	address string
	timeout time.Duration
	result  Result
	errs    []error
	lock    *semaphore.Weighted
}

func newUDPOpts(address string, ports []int, timeout time.Duration) (udpOpts, error) {

	scan := udpOpts{address: address, timeout: timeout}

	for i := range ports {
		scan.result = append(scan.result, Port{Port: ports[i], State: FILTERED})
	}

	ulimit, err := GetMaxFD()
	if err != nil {
		return scan, fmt.Errorf("failed to get max file descriptors: %s", err)
	}
	// With ulimit/2, scan will bit slower, but dont have to bother with " socket: too many open files" errors.
	scan.lock = semaphore.NewWeighted(int64(ulimit) / 2)

	return scan, nil
}

func (c *udpOpts) connect(p *Port, wg *sync.WaitGroup, e chan<- error) {

	defer c.lock.Release(1)
	defer wg.Done()
	c.lock.Acquire(context.TODO(), 1)

	target := fmt.Sprintf("%s:%d", c.address, p.Port)

	conn, err := net.DialTimeout("udp", target, c.timeout)
	if err != nil {
		e <- fmt.Errorf("dial error: %s", err)
		return
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		e <- fmt.Errorf("failed to set deadline: %s", err)
		return
	}

	// Write some random string
	buf := make([]byte, rand.Intn(64))

	if _, err = rand.Read(buf); err != nil {
		e <- fmt.Errorf("failed to read random: %s", err)
		return
	}

	_, err = conn.Write([]byte("elmasy.com\r\n"))
	if err != nil {

		switch true {
		case strings.Contains(err.Error(), "i/o timeout"):
			// Port filtered or open, report as FILTERED
			return
		default:
			e <- fmt.Errorf("unknown write error: %s", err)
			return
		}
	}

	// Try to read something
	n, err := conn.Read(buf)
	if err != nil {

		switch true {
		case strings.Contains(err.Error(), "i/o timeout"):
			// Port filtered or open, report as FILTERED
			return
		case strings.Contains(err.Error(), "connection refused"):
			// Cleraly closed
			p.State = CLOSED
			return
		default:
			e <- fmt.Errorf("unknown write error: %s", err)
			return
		}
	}

	if n > 0 {
		// Read more than one byte, clearly OPEN
		p.State = OPEN
	}
}

func udpscan(opts udpOpts) (Result, []error) {

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

func UDPScan(address string, ports []int, timeout time.Duration) (Result, []error) {

	var results Result

	scan1, err := newUDPOpts(address, ports, timeout)
	if err != nil {
		return nil, []error{err}
	}

	scan1r, errs := udpscan(scan1)
	if err != nil {
		return nil, errs
	}

	results = append(results, scan1r.GetPorts(OPEN)...)
	results = append(results, scan1r.GetPorts(CLOSED)...)

	if scan1r.Len(FILTERED) > 0 {

		scan2, err := newUDPOpts(address, scan1r.GetPortsInt(FILTERED), timeout)
		if err != nil {
			return nil, []error{err}
		}

		scan2r, errs := udpscan(scan2)
		if errs != nil {
			return nil, errs
		}

		results = append(results, scan2r...)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Port < results[j].Port })

	return results, nil
}
