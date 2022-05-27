package sdk

import (
	"fmt"
	"io"
	"net/http"
)

var (
	API_PATH      = "https://elmasy.com/api"
	USER_AGENT    = "Elmasy-SDK"
	DefaultClient = &http.Client{}
)

// Get query API_PATH + endpoint.
// Returns the body, the status code and error.
func Get(endpoint string) ([]byte, int, error) {

	req, err := http.NewRequest("GET", API_PATH+endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create a new request: %s", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", USER_AGENT)

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return body, resp.StatusCode, err
}
