package sdk

import (
	"encoding/json"
	"fmt"
)

func GetIP() (string, error) {

	body, status, err := Get("/ip")
	if err != nil {
		return "", err
	}

	switch status {
	case 200:
		result := Result{}

		if err := json.Unmarshal(body, &result); err != nil {
			return "", fmt.Errorf("failed to unmarshal: %s", err)
		}
		return result.Result, nil
	default:
		return "", fmt.Errorf("unknown status: %d", status)
	}

}
