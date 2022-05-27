package sdk

type ResultStr struct {
	Result string `json:"result"`
}

type ResultBool struct {
	Result bool `json:"result"`
}

type ResultStrs struct {
	Results []string `json:"results"`
}

type Error struct {
	Err string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

type Cipher struct {
	Name     string `json:"name"`
	Security string `json:"security"`
}

type PubKey struct {
	Algo     string `json:"algo"`
	Size     int    `json:"size"` // Bits
	Key      string `json:"key"`
	Exponent int    `json:"exponent"`
	Modulus  string `json:"modulus"`
}

// Additional is the additional certificates (eg.: intermediate cert)
type Additional struct {
	CommonName         string `json:"commonname"`
	Hash               string `json:"hash"`
	NotAfter           string `json:"notafter"`
	Issuer             string `json:"issuer"`
	PublicKey          PubKey `json:"publickey"`
	SignatureAlgorithm string `json:"signaturealgorithm"`
}

// Cert is hold the fields "interesting" part of the certficate chain.
type Cert struct {
	CommonName         string       `json:"commonname"`
	Hash               string       `json:"hash"` // SHA256
	AlternativeNames   []string     `json:"alternativenames"`
	SignatureAlgorithm string       `json:"signaturealgorithm"`
	PublicKey          PubKey       `json:"publickey"`
	SerialNumber       string       `json:"serialnumber"`
	Issuer             string       `json:"issuer"`
	NotBefore          string       `json:"notbefore"`
	NotAfter           string       `json:"notafter"`
	Verified           bool         `json:"verified"`
	VerifiedError      string       `json:"verifiederror"`
	Chain              []Additional `json:"chain"`
}

type TLSVersion struct {
	Version   string   `json:"version"`
	Supported bool     `json:"supported"`
	Ciphers   []Cipher `json:"ciphers"`
}

type TLS struct {
	Versions []TLSVersion `json:"version"`
	Cert     Cert         `json:"cert"`
}

type Target struct {
	IP  string `json:"ip"`
	TLS TLS    `json:"tls"`
}

type Result struct {
	Domain  string   `json:"domain"`
	Targets []Target `json:"targets"`
	Errors  []string `json:"errors"`
}
