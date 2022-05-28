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
	CommonName         string `json:"commonName"`
	Hash               string `json:"hash"`
	NotAfter           string `json:"notAfter"`
	Issuer             string `json:"issuer"`
	PublicKey          PubKey `json:"publicKey"`
	SignatureAlgorithm string `json:"signatureAlgorithm"`
}

// Cert is hold the fields "interesting" part of the certficate chain.
type Cert struct {
	CommonName         string       `json:"commonName"`
	Hash               string       `json:"hash"` // SHA256
	AlternativeNames   []string     `json:"alternativeNames"`
	SignatureAlgorithm string       `json:"signatureAlgorithm"`
	PublicKey          PubKey       `json:"publicKey"`
	SerialNumber       string       `json:"serialNumber"`
	Issuer             string       `json:"issuer"`
	NotBefore          string       `json:"notBefore"`
	NotAfter           string       `json:"notAfter"`
	Verified           bool         `json:"verified"`
	VerifiedError      string       `json:"verifiedError"`
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
