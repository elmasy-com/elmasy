package sdk

import (
	"encoding/json"
	"fmt"
)

func Scan(target, port, network string) (Result, error) {

	url := fmt.Sprintf("/scan?target=%s&port=%s&network=%s", target, port, network)

	body, status, err := Get(url)
	if err != nil {
		return Result{}, err
	}

	switch status {
	case 200:
		var r Result

		if err := json.Unmarshal(body, &r); err != nil {
			return r, fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r, nil
	case 400, 403, 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return Result{}, fmt.Errorf("failed to unmarshal: %s", err)
		}
		return Result{}, e
	default:
		return Result{}, fmt.Errorf("unknown status: %d", status)
	}

}

func PortScan(technique, ip, ports, timeout string) (string, error) {

	url := fmt.Sprintf("/scan/port?technique=%s&ip=%s&port=%s&timeout=%s", technique, ip, ports, timeout)

	body, status, err := Get(url)
	if err != nil {
		return "", err
	}

	switch status {
	case 200:
		var r ResultStr

		if err := json.Unmarshal(body, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}

		return r.Result, nil
	case 400, 403:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}

		return "", e
	case 500:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}

		return "", e
	default:
		return "", fmt.Errorf("unknown status: %d", status)
	}
}
