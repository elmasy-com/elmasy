package sdk

import (
	"encoding/json"
	"fmt"
)

func Scan(target, port, network string) ([]Target, error) {

	url := fmt.Sprintf("/scan?target=%s&port=%s&network=%s", target, port, network)

	body, status, err := Get(url)
	if err != nil {
		return nil, err
	}

	switch status {
	case 200:
		var r []Target

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

func PortScan(technique, ip, ports string) (Ports, []error) {

	url := fmt.Sprintf("/scan/port?technique=%s&ip=%s&ports=%s", technique, ip, ports)

	body, status, err := Get(url)
	if err != nil {
		return nil, []error{err}
	}

	switch status {
	case 200:
		var r Ports

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, []error{fmt.Errorf("failed to unmarshal: %s", err)}
		}

		return r, nil
	case 400, 403:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return nil, []error{fmt.Errorf("failed to unmarshal: %s", err)}
		}

		return nil, []error{e}
	case 500:
		e := Errors{}

		if err := json.Unmarshal(body, &e); err != nil {
			return nil, []error{fmt.Errorf("failed to unmarshal: %s", err)}
		}

		errs := make([]error, 0)

		for i := range e.Errors {
			errs = append(errs, fmt.Errorf(e.Errors[i]))
		}

		return nil, errs
	default:
		return nil, []error{fmt.Errorf("unknown status: %d", status)}
	}
}

func IsPortOpen(technique, ip, port string) (bool, error) {

	ports, errs := PortScan(technique, ip, port)
	if errs != nil {
		return false, fmt.Errorf("%v", errs)
	}

	if len(ports) != 1 {
		return false, fmt.Errorf("multiple ports: %v", ports)
	}

	return ports[0].State == "open", nil
}
