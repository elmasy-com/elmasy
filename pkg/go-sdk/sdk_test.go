package sdk

import "testing"

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Log(ip)
}

func TestGetRandomIP(t *testing.T) {

	rand, err := GetRandomIP("4")
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Logf("%s", rand)
}

func TestGetRandomPort(t *testing.T) {

	port, err := GetRandomPort()
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Logf("%s", port)
}

func TestDNSLookup(t *testing.T) {

	r, err := DNSLookup("A", "elmasy.com")
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Logf("%v", r)
}

func TestAnalyzeTLS(t *testing.T) {

	r, err := AnalyzeTLS("tls12", "tcp", "95.216.184.245", "443")
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Logf("%v", r)
}

func TestPortScan(t *testing.T) {

	r, errs := PortScan("stealth", "95.216.184.245", "80,443")
	if errs != nil {
		t.Fatalf("FAIL: %v", errs)
	}

	t.Logf("%v", r)
}

func TestProbe(t *testing.T) {

	r, err := Probe("tls12", "tcp", "95.216.184.245", "443")
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Logf("%v", r)
}

func TestScan(t *testing.T) {

	r, err := Scan("elmasy.com", "443", "tcp")
	if err != nil {
		t.Fatalf("FAIL: %s", err)
	}

	t.Logf("%#v", r)
}
