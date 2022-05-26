package types

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
	Err string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

type Cipher struct {
	Name     string `json:"name"`
	Security string `json:"security"`
}

type TLS struct {
	IP        string   `json:"ip"`
	Version   string   `json:"version"`
	Supported bool     `json:"supported"`
	Ciphers   []Cipher `json:"ciphers"`
}

type Target struct {
	Target string `json:"target"`
	TLS    []TLS  `json:"tls"`
}

type Port struct {
	Port  string
	State string
}

type Ports []Port

func (e Error) Error() string {
	return e.Err
}
