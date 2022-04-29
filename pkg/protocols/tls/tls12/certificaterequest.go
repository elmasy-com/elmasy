package tls12

import (
	"fmt"

	"github.com/elmasy-com/bytebuilder"
)

/*
   enum {
       rsa_sign(1), dss_sign(2), rsa_fixed_dh(3), dss_fixed_dh(4),
       rsa_ephemeral_dh(5), dss_ephemeral_dh(6), fortezza_kea(20),
       (255)
   } ClientCertificateType;

   opaque DistinguishedName<1..2^16-1>;

   struct {
       ClientCertificateType certificate_types<1..2^8-1>;
       DistinguishedName certificate_authorities<3..2^16-1>;
   } CertificateRequest;
*/

type certificateRequest struct {
	CertTypes []uint8
	CertAuths []byte
}

func unmarshalCertificateRequest(bytes []byte) (certificateRequest, error) {

	var (
		csr  = certificateRequest{}
		buf  = bytebuilder.NewBuffer(bytes)
		tLen uint8
		ok   bool
	)

	if tLen, ok = buf.ReadUint8(); !ok {
		return csr, fmt.Errorf("failed to read certifcate_types length")
	}

	for i := 0; i < int(tLen); i++ {
		var t uint8

		if t, ok = buf.ReadUint8(); !ok {
			return csr, fmt.Errorf("failed to read certificate type")
		}

		csr.CertTypes = append(csr.CertTypes, t)
	}

	if csr.CertAuths, ok = buf.ReadVector(16); !ok {
		return csr, fmt.Errorf("failed to read certificate_authorities")
	}

	return csr, nil

}
