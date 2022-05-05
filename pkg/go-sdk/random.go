package sdk

import (
	"encoding/json"
	"fmt"
)

func GetRandomIP(version string) (string, error) {

	body, status, err := Get("/random/ip/" + version)
	if err != nil {
		return "", err
	}

	switch status {
	case 200:
		r := Result{}

		if err := json.Unmarshal(body, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return r.Result, nil
	case 400:
		e := Error{}

		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return "", e
	default:
		return "", fmt.Errorf("unknown status: %d", status)
	}
}

func GetRandomPort() (string, error) {

	body, status, err := Get("/random/port")
	if err != nil {
		return "", err
	}

	switch status {
	case 200:
		r := Result{}

		if err := json.Unmarshal(body, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return r.Result, nil
	default:
		return "", fmt.Errorf("unknown status: %d", status)
	}

}
