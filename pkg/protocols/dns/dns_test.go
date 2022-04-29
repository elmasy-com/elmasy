package dns

import (
	"testing"
	"time"
)

func TestQueryA(t *testing.T) {

	MAX_RETRIES = 10

	_, err := QueryA("elmasy.com")
	if err != nil {
		t.Fatalf("TestQueryA failed: %s\n", err)
	}

}

func TestQueryAAAA(t *testing.T) {

	MAX_RETRIES = 10

	_, err := QueryAAAA("elmasy.com")
	if err != nil {
		t.Fatalf("TestQueryAAAA failed: %s\n", err)
	}
}

func TestQueryMX(t *testing.T) {

	MAX_RETRIES = 10

	_, err := QueryMX("elmasy.com")
	if err != nil {
		t.Fatalf("TestQueryMX failed: %s\n", err)
	}
}

func TestQueryTXT(t *testing.T) {

	MAX_RETRIES = 10

	_, err := QueryTXT("elmasy.com")
	if err != nil {
		t.Fatalf("TestQueryTXT failed: %s\n", err)
	}
}

func TestProbe(t *testing.T) {

	e, err := Probe("udp", "1.1.1.1", "53", 2*time.Second)
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	if !e {
		t.Fatalf("TestProbe failed: 1.1.1.1:53 should be a valid DNS server")
	}
}
