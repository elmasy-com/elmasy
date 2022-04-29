package dns

import "testing"

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
