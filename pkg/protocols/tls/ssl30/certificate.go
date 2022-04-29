package ssl30

import (
	"crypto/x509"
	"fmt"

	"github.com/elmasy-com/bytebuilder"
)

/*
	opaque ASN.1Cert<1..2^24-1>;
	struct {
    	ASN.1Cert certificate_list<1..2^24-1>;
	} Certificate;
*/

type certificate struct {
	Certificates []x509.Certificate
}

func unmarhsalCertificate(bytes []byte) (certificate, error) {

	var (
		cert certificate
		ok   bool
		buf  = bytebuilder.NewBuffer(bytes)
	)

	if ok = buf.Skip(3); !ok {
		return cert, fmt.Errorf("failed to skip certificate_list length")
	}

	for !buf.Empty() {

		certbytes, ok := buf.ReadVector(24)
		if !ok {
			return cert, fmt.Errorf("failed to read certificate vector")
		}

		cer, err := x509.ParseCertificate(certbytes)
		if err != nil {
			return cert, fmt.Errorf("failed to parse certificate: %s", err)
		}

		cert.Certificates = append(cert.Certificates, *cer)
	}

	return cert, nil
}
