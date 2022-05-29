package certificate

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"

	etls "github.com/elmasy-com/elmasy/pkg/protocols/tls"
	"golang.org/x/crypto/ocsp"
)

type PubKey struct {
	Algo x509.PublicKeyAlgorithm
	Key  any // *rsa.PublicKey, *ed25519.PublicKey, ...
}

// Additional is the additional certificates (eg.: intermediate cert)
type Additional struct {
	CommonName         string
	Hash               [32]byte
	NotAfter           time.Time
	Issuer             string
	PublicKey          PubKey
	SignatureAlgorithm x509.SignatureAlgorithm
}

// Cert is hold the fields "interesting" part of the certficate chain.
type Cert struct {
	CommonName         string
	Hash               [32]byte // SHA256
	AlternativeNames   []string
	SignatureAlgorithm x509.SignatureAlgorithm
	PublicKey          PubKey
	SerialNumber       *big.Int
	Issuer             string
	NotBefore          time.Time
	NotAfter           time.Time
	Verified           bool
	VerifiedError      error // This is set if Verified == false
	Chain              []Additional
}

// Ordered by usage
var tlsVersions = []string{"tls12", "tls13", "tls11", "tls10", "ssl30"}

func Grab(network, ip, port string, timeout time.Duration, servername string) ([]x509.Certificate, error) {

	for i := range tlsVersions {
		r, err := etls.Handshake(tlsVersions[i], network, ip, port, timeout, servername)
		if err != nil {
			return nil, err
		}

		if r.Supported {
			return r.Certificates, nil
		}
	}

	return nil, fmt.Errorf("TLS not supported")
}

func verifyOCSP(leaf x509.Certificate, issuer x509.Certificate) error {

	opts := ocsp.RequestOptions{Hash: crypto.SHA1}

	buf, err := ocsp.CreateRequest(&leaf, &issuer, &opts)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, leaf.OCSPServer[0], bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	url, err := url.Parse(leaf.OCSPServer[0])
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/ocsp-request")
	req.Header.Add("Accept", "application/ocsp-response")
	req.Header.Add("host", url.Host)

	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ocspResp, err := ocsp.ParseResponseForCert(body, &leaf, &issuer)
	if err != nil {
		return err
	}

	switch ocspResp.Status {
	case ocsp.Good:
		return nil
	case ocsp.Revoked:
		return fmt.Errorf("revoked")
	case ocsp.Unknown:
		return fmt.Errorf("unknown")
	default:
		return fmt.Errorf("invalid status: %d", ocspResp.Status)
	}
}

func Verify(certs []x509.Certificate, servername string) error {

	if len(certs) < 1 {
		return fmt.Errorf("zero certificate given")
	}

	opts := x509.VerifyOptions{DNSName: servername}

	if len(certs) > 1 {
		pool := x509.NewCertPool()
		chain := certs[1:]
		for i := range chain {
			pool.AddCert(&chain[i])
		}
		opts.Intermediates = pool
	}

	if _, err := certs[0].Verify(opts); err != nil {
		// remove the "x509: " prefix
		e := strings.TrimPrefix(err.Error(), "x509: ")
		return fmt.Errorf("%s", e)
	}

	if len(certs[0].OCSPServer) > 0 && len(certs) >= 2 {

		if err := verifyOCSP(certs[0], certs[1]); err != nil {
			return err
		}
	}

	return nil
}

// parseLeafCert parse the leaf cert and fill the fields of r (result Cert)
func parseLeafCert(cert x509.Certificate, r *Cert) {

	r.CommonName = cert.Subject.CommonName

	r.Hash = sha256.Sum256(cert.Raw)

	for i := range cert.DNSNames {
		r.AlternativeNames = append(r.AlternativeNames, cert.DNSNames[i])
	}
	for i := range cert.IPAddresses {
		r.AlternativeNames = append(r.AlternativeNames, cert.IPAddresses[i].String())
	}

	r.SignatureAlgorithm = cert.SignatureAlgorithm
	r.PublicKey.Algo = cert.PublicKeyAlgorithm
	r.PublicKey.Key = cert.PublicKey

	r.SerialNumber = cert.SerialNumber

	r.Issuer = cert.Issuer.CommonName

	r.NotBefore = cert.NotBefore
	r.NotAfter = cert.NotAfter
}

func parseChain(certs []x509.Certificate, r *Cert) {

	for i := range certs {
		a := Additional{}

		a.CommonName = certs[i].Subject.CommonName
		a.Hash = sha256.Sum256(certs[i].Raw)
		a.NotAfter = certs[i].NotAfter
		a.Issuer = certs[i].Issuer.CommonName
		a.PublicKey.Algo = certs[i].PublicKeyAlgorithm
		a.PublicKey.Key = certs[i].PublicKey
		a.SignatureAlgorithm = certs[i].SignatureAlgorithm

		r.Chain = append(r.Chain, a)
	}
}

func Scan(network, ip, port string, timeout time.Duration, servername string) (Cert, error) {

	if servername == "" {
		return Cert{}, fmt.Errorf("servername is empty")
	}

	result := Cert{}

	certs, err := Grab(network, ip, port, timeout, servername)
	if err != nil {
		return Cert{}, err
	}

	if len(certs) == 0 {
		return result, fmt.Errorf("no certificate")
	}

	if err := Verify(certs, servername); err == nil {
		result.Verified = true
	} else {
		result.VerifiedError = err
	}

	parseLeafCert(certs[0], &result)

	if len(certs) > 1 {
		parseChain(certs[1:], &result)
	}

	return result, nil
}
