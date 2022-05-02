package sdk

import (
	"encoding/json"
	"fmt"
	"io"
)

func DNSLookup(t, n string) ([]string, error) {

	url := fmt.Sprintf("/protocol/dns/%s/%s", t, n)

	resp, err := Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err)
	}

	switch resp.StatusCode {
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
		return nil, fmt.Errorf("unknown status: %s", resp.Status)
	}
}

func AnalyzeTLS(version, network, ip, port string) (TLS, error) {

	url := fmt.Sprintf("/protocol/tls?version=%s&network=%s&ip=%s&port=%s", version, network, ip, port)

	resp, err := Get(url)
	if err != nil {
		return TLS{}, fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TLS{}, fmt.Errorf("failed to read body: %s", err)
	}

	switch resp.StatusCode {
	case 200:
		r := TLS{}

		if err := json.Unmarshal(body, &r); err != nil {
			return TLS{}, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return TLS{}, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return TLS{}, e
	default:
		return TLS{}, fmt.Errorf("unknown status: %s", resp.Status)
	}
}

func Probe(protocol, network, ip, port string) (bool, error) {

	url := fmt.Sprintf("/protocol/probe?protocol=%s&network=%s&ip=%s&port=%s", protocol, network, ip, port)

	resp, err := Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read body: %s", err)
	}

	switch resp.StatusCode {
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
		return false, fmt.Errorf("unknown status: %s", resp.Status)
	}
}
