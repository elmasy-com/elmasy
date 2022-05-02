package sdk

import (
	"encoding/json"
	"fmt"
	"io"
)

func GetIP() (string, error) {

	resp, err := Get("/ip")
	if err != nil {
		return "", fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %s", err)
	}

	switch resp.StatusCode {
	case 200:
		result := Result{}

		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return result.Result, nil
	default:
		return "", fmt.Errorf("unknown status: %s", resp.Status)
	}

}
