package sdk

import (
	"encoding/json"
	"fmt"
	"io"
)

func GetRandomIP(version string) (string, error) {

	resp, err := Get("/random/ip/" + version)
	if err != nil {
		return "", fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %s", err)
	}

	r := Result{}
	e := Error{}

	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(body, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return r.Result, nil
	case 400:
		if err := json.Unmarshal(body, &e); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return "", e
	default:
		return "", fmt.Errorf("unknown status: %s", resp.Status)
	}
}

func GetRandomPort() (string, error) {

	resp, err := Get("/random/port")
	if err != nil {
		return "", fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %s", err)
	}

	r := Result{}

	switch resp.StatusCode {
	case 200:
		if err := json.Unmarshal(body, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return r.Result, nil
	default:
		return "", fmt.Errorf("unknown status: %s", resp.Status)
	}

}
