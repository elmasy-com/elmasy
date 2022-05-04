package sdk

import (
	"encoding/json"
	"fmt"
	"io"
)

func PortScan(technique, ip, ports string) ([]Ports, []error) {

	url := fmt.Sprintf("/scan/port?technique=%s&ip=%s&ports=%s", technique, ip, ports)

	resp, err := Get(url)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to query: %s", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read body: %s", err)}
	}

	switch resp.StatusCode {
	case 200:
		r := struct {
			Result []Ports `json:"result"`
		}{}

		if err := json.Unmarshal(body, &r); err != nil {
			return nil, []error{fmt.Errorf("failed to unmarshal: %s", err)}
		}

		return r.Result, nil
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
		return nil, []error{fmt.Errorf("unknown status: %s", resp.Status)}
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
