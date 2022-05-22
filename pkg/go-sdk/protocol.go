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
		r := Results{}

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

func AnalyzeTLS(version, network, ip, port, servername string) ([]TLS, error) {

	url := fmt.Sprintf("/protocol/tls?version=%s&network=%s&target=%s&port=%s&servername=%s", version, network, ip, port, servername)

	body, status, err := Get(url)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		var r []TLS

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return nil, e
	default:
		return nil, fmt.Errorf("unknown status: %d", status)
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
		r := Result{}

		if err = json.Unmarshal(body, &r); err != nil {
			return false, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Result == "true", nil
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
