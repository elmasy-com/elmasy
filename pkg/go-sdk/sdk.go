package sdk

import (
	"fmt"
	"net/http"
)

var (
	API_PATH                   = "https://elmasy.com/api"
	DefaultClient *http.Client = &http.Client{}
)

type Result struct {
	Result string `json:"result"`
}

type ResultBool struct {
	Result bool `json:"result"`
}

type Results struct {
	Results []string `json:"results"`
}

type Error struct {
	E string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

type TLS struct {
	Supported bool     `json:"supported"`
	Ciphers   []string `json:"ciphers"`
}

type Ports struct {
	Port  string
	State string
}

func (e Error) Error() string {
	return e.E
}

func Get(endpoint string) (*http.Response, error) {

	req, err := http.NewRequest("GET", API_PATH+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %s", err)
	}
	req.Header.Add("Accept", "application/json")

	return DefaultClient.Do(req)
}
