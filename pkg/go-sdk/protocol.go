package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/elmasy-com/elmasy/pkg/types"
)

func DNSLookup(t, n string) ([]string, error) {

	url := fmt.Sprintf("/protocol/dns/%s/%s", t, n)

	body, status, err := Get(url)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		r := types.Results{}

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Results, nil
	case 400, 404, 500:
		e := types.Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return nil, e
	default:
		return nil, fmt.Errorf("unknown status: %d", status)
	}
}

func AnalyzeTLS(version, network, ip, port, servername string) ([]types.TLS, error) {

	url := fmt.Sprintf("/protocol/tls?version=%s&network=%s&target=%s&port=%s&servername=%s", version, network, ip, port, servername)

	body, status, err := Get(url)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		var r []types.TLS

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 500:
		e := types.Error{}

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
		r := types.ResultBool{}

		if err = json.Unmarshal(body, &r); err != nil {
			return false, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Result, nil
	case 400, 500:
		e := types.Error{}

		if err = json.Unmarshal(body, &e); err != nil {
			return false, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return false, e
	default:
		return false, fmt.Errorf("unknown status: %d", status)
	}
}
