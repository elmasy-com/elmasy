package sdk

import (
	"encoding/json"
	"fmt"
)

func DNSLookup(t, n string) ([]string, error) {

	url := fmt.Sprintf("/protocol/dns/%s/%s", t, n)

	body, status, err := Get(url)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		r := ResultStrs{}

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Results, nil
	case 400, 404, 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return nil, e
	default:
		return nil, fmt.Errorf("unknown status: %d", status)
	}
}

func AnalyzeTLS(version, network, ip, port, servername string) (TLSVersion, error) {

	url := fmt.Sprintf("/protocol/tls?version=%s&network=%s&ip=%s&port=%s&servername=%s", version, network, ip, port, servername)

	body, status, err := Get(url)
	if err != nil {
		return TLSVersion{}, err
	}

	switch status {
	case 200:
		var r TLSVersion

		if err := json.Unmarshal(body, &r); err != nil {
			return TLSVersion{}, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return TLSVersion{}, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return TLSVersion{}, e
	default:
		return TLSVersion{}, fmt.Errorf("unknown status: %d", status)
	}
}

func GetCertificate(network, ip, port, servername string) (Cert, error) {

	url := fmt.Sprintf("/protocol/tls/certificate?network=%s&ip=%s&port=%s&servername=%s", network, ip, port, servername)

	body, status, err := Get(url)
	if err != nil {
		return Cert{}, err
	}

	switch status {
	case 200:
		var r Cert

		if err := json.Unmarshal(body, &r); err != nil {
			return Cert{}, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return Cert{}, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return Cert{}, e
	default:
		return Cert{}, fmt.Errorf("unknown status: %d", status)
	}

}

func Probe(protocol, network, ip, port string) (bool, error) {

	url := fmt.Sprintf("/protocol/probe?protocol=%s&network=%s&ip=%s&port=%s", protocol, network, ip, port)

	body, status, err := Get(url)
	if err != nil {
		return false, err
	}

	switch status {
	case 200:
		r := ResultBool{}

		if err = json.Unmarshal(body, &r); err != nil {
			return false, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Result, nil
	case 400, 500:
		e := Error{}

		if err = json.Unmarshal(body, &e); err != nil {
			return false, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return false, e
	default:
		return false, fmt.Errorf("unknown status: %d", status)
	}
}
