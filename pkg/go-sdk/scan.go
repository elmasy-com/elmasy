package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/elmasy-com/elmasy/pkg/types"
)

func Scan(target, port, network string) ([]types.Target, error) {

	url := fmt.Sprintf("/scan?target=%s&port=%s&network=%s", target, port, network)

	body, status, err := Get(url)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		var r []types.Target

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 403, 500:
		e := types.Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %s", err)
		}
		return nil, e
	default:
		return nil, fmt.Errorf("unknown status: %d", status)
	}

}

func PortScan(technique, ip, ports, timeout string) (string, error) {

	url := fmt.Sprintf("/scan/port?technique=%s&ip=%s&port=%s&timeout=%s", technique, ip, ports, timeout)

	body, status, err := Get(url)
	if err != nil {
		return "", err
	}

	fmt.Printf("BODY:\n%s\n\n", body)
	switch status {
	case 200:
		var r types.Result

		if err := json.Unmarshal(body, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Result, nil
	case 400, 403:
		e := types.Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}

		return "", e
	case 500:
		e := types.Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}

		return "", e
	default:
		return "", fmt.Errorf("unknown status: %d", status)
	}
}
