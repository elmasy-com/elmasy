package certificate

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/elmasy-com/elmasy/internal/utils"
	"github.com/elmasy-com/elmasy/pkg/go-sdk"
	ecertificate "github.com/elmasy-com/elmasy/pkg/protocols/tls/certificate"
	"github.com/elmasy-com/identify"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {

	network := c.DefaultQuery("network", "tcp")
	if network != "tcp" && network != "udp" {
		err := fmt.Errorf("Invalid network: %s", network)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	ip := c.Query("ip")
	if ip == "" {
		err := fmt.Errorf("ip is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}
	if !identify.IsValidIP(ip) {
		err := fmt.Errorf("Invalid ip: %s", ip)
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}
	ip = utils.IPv6BracketAdd(ip)

	port := c.Query("port")
	if port == "" {
		err := fmt.Errorf("port is empty")
		c.Error(err)
		c.JSON(http.StatusBadRequest, sdk.Error{Err: err.Error()})
		return
	}

	servername := c.Query("servername")

	r, err := ecertificate.Scan(network, ip, port, 2*time.Second, servername)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, sdk.Error{Err: err.Error()})
		return
	}

	cert := sdk.Cert{}

	cert.CommonName = r.CommonName
	cert.Hash = fmt.Sprintf("%x", r.Hash)
	cert.AlternativeNames = r.AlternativeNames
	cert.SignatureAlgorithm = r.SignatureAlgorithm.String()
	cert.PublicKey.Algo = r.PublicKey.Algo.String()

	switch v := r.PublicKey.Key.(type) {
	case *rsa.PublicKey:
		cert.PublicKey.Exponent = v.E
		cert.PublicKey.Modulus = fmt.Sprintf("%x", v.N)
		cert.PublicKey.Size = v.Size() * 8
	case *dsa.PublicKey:
		cert.PublicKey.Key = fmt.Sprintf("%x", v.Y)
		cert.PublicKey.Size = v.P.BitLen()
	case *ecdsa.PublicKey:
		cert.PublicKey.Key = fmt.Sprintf("%x", elliptic.Marshal(v.Curve, v.X, v.Y))
		cert.PublicKey.Size = v.Params().BitSize
	case *ed25519.PublicKey:
		cert.PublicKey.Key = fmt.Sprintf("%x", *v)
		cert.PublicKey.Size = len(*v) * 16
	default:
		err := fmt.Errorf("unknown public key")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, sdk.Error{Err: err.Error()})
		return
	}

	cert.SerialNumber = fmt.Sprintf("%x", r.SerialNumber)
	cert.Issuer = r.Issuer
	cert.NotBefore = r.NotBefore.Format(time.RFC1123)
	cert.NotAfter = r.NotAfter.Format(time.RFC1123)
	cert.Verified = r.Verified
	if r.VerifiedError != nil {
		cert.VerifiedError = r.VerifiedError.Error()
	}

	for i := range r.Chain {
		a := sdk.Additional{}
		a.CommonName = r.Chain[i].CommonName
		a.Hash = fmt.Sprintf("%x", r.Chain[i].Hash)
		a.NotAfter = r.Chain[i].NotAfter.Format(time.RFC1123)
		a.Issuer = r.Chain[i].Issuer
		a.PublicKey.Algo = r.Chain[i].PublicKey.Algo.String()

		switch v := r.Chain[i].PublicKey.Key.(type) {
		case *rsa.PublicKey:
			a.PublicKey.Exponent = v.E
			a.PublicKey.Modulus = fmt.Sprintf("%x", v.N)
			a.PublicKey.Size = v.Size() * 8
		case *dsa.PublicKey:
			a.PublicKey.Key = fmt.Sprintf("%x", v.Y)
			a.PublicKey.Size = v.P.BitLen()
		case *ecdsa.PublicKey:
			a.PublicKey.Key = fmt.Sprintf("%x", elliptic.Marshal(v.Curve, v.X, v.Y))
			a.PublicKey.Size = v.Params().BitSize
		case *ed25519.PublicKey:
			a.PublicKey.Key = fmt.Sprintf("%x", *v)
			a.PublicKey.Size = len(*v) * 16
		default:
			err := fmt.Errorf("unknown public key in chain")
			c.Error(err)
			c.JSON(http.StatusInternalServerError, sdk.Error{Err: err.Error()})
			return
		}

		a.SignatureAlgorithm = r.Chain[i].SignatureAlgorithm.String()

		cert.Chain = append(cert.Chain, a)
	}

	c.JSON(http.StatusOK, cert)
}
