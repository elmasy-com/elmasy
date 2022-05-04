package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/elmasy-com/elmasy/pkg/protocols/tls/ciphersuite"
)

/*
	Parse every ciphersuite from ciphersuite.info and print it to stdout formatted as a Go struct.
*/
type CipherInfo struct {
	ByteOne  string   `json:"hex_byte_1"`
	ByteTwo  string   `json:"hex_byte_2"`
	Versions []string `json:"tls_version"`
	Security string   `json:"security"`
}

type Cipher struct {
	Cipher map[string]interface{}
}

type Ciphers struct {
	Ciphers []map[string]CipherInfo `json:"ciphersuites"`
}

var CipherSuites []ciphersuite.CipherSuite

func main() {

	resp, err := http.Get("https://ciphersuite.info/api/cs/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to GET ciphersuite.info API: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read body: %s\n", err)
		os.Exit(1)
	}

	cs := Ciphers{}

	if err := json.Unmarshal(body, &cs); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("package ciphersuite\n\n")
	fmt.Printf("var CipherSuites = []CipherSuite {\n")

	for i := range cs.Ciphers {
		printCipherSuites(cs.Ciphers[i])
	}

	// This 2 ciphersuite not exist on ciphersuite.info, print is manually.
	fmt.Printf("{[]byte{0x00, 0x1C}, \"SSL_FORTEZZA_KEA_WITH_NULL_SHA\", []uint16{SSL30}, \"weak\"},\n")
	fmt.Printf("{[]byte{0x00, 0x1D}, \"SSL_FORTEZZA_KEA_WITH_FORTEZZA_CBC_SHA\", []uint16{SSL30}, \"weak\"},\n")

	fmt.Printf("}\n")

}

func printCipherSuites(m map[string]CipherInfo) {

	for k, v := range m {

		fmt.Printf("{[]byte{%s, %s}, ", v.ByteOne, v.ByteTwo)
		fmt.Printf("\"%s\", ", k)

		vers := make([]string, 0)
		if isSSL30Cipher(k) {
			vers = append(vers, "SSL30")
		}
		for _, version := range v.Versions {
			switch version {
			case "TLS1.0":
				vers = append(vers, "TLS10")
			case "TLS1.1":
				vers = append(vers, "TLS11")
			case "TLS1.2":
				vers = append(vers, "TLS12")
			case "TLS1.3":
				vers = append(vers, "TLS13")

			}
		}

		fmt.Printf("[]uint16{%s}, \"%s\"},\n", strings.Join(vers, ","), v.Security)

	}
}

var ssl30ciphers = []string{

	"SSL_FORTEZZA_KEA_WITH_NULL_SHA",
	"SSL_FORTEZZA_KEA_WITH_FORTEZZA_CBC_SHA",

	"TLS_NULL_WITH_NULL_NULL",
	"TLS_RSA_WITH_NULL_MD5",
	"TLS_RSA_WITH_NULL_SHA",
	"TLS_RSA_EXPORT_WITH_RC4_40_MD5",
	"TLS_RSA_WITH_RC4_128_MD5",
	"TLS_RSA_WITH_RC4_128_SHA",
	"TLS_RSA_EXPORT_WITH_RC2_CBC_40_MD5",
	"TLS_RSA_WITH_IDEA_CBC_SHA",
	"TLS_RSA_EXPORT_WITH_DES40_CBC_SHA",
	"TLS_RSA_WITH_DES_CBC_SHA",
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_DH_DSS_EXPORT_WITH_DES40_CBC_SHA",
	"TLS_DH_DSS_WITH_DES_CBC_SHA",
	"TLS_DH_DSS_WITH_3DES_EDE_CBC_SHA",
	"TLS_DH_RSA_EXPORT_WITH_DES40_CBC_SHA",
	"TLS_DH_RSA_WITH_DES_CBC_SHA",
	"TLS_DH_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_DHE_DSS_EXPORT_WITH_DES40_CBC_SHA",
	"TLS_DHE_DSS_WITH_DES_CBC_SHA",
	"TLS_DHE_DSS_WITH_3DES_EDE_CBC_SHA",
	"TLS_DHE_RSA_EXPORT_WITH_DES40_CBC_SHA",
	"TLS_DHE_RSA_WITH_DES_CBC_SHA",
	"TLS_DHE_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_DH_anon_EXPORT_WITH_RC4_40_MD5",
	"TLS_DH_anon_WITH_RC4_128_MD5",
	"TLS_DH_anon_EXPORT_WITH_DES40_CBC_SHA",
	"TLS_DH_anon_WITH_DES_CBC_SHA",
	"TLS_DH_anon_WITH_3DES_EDE_CBC_SHA",

	"TLS_KRB5_WITH_DES_CBC_SHA",
}

func isSSL30Cipher(name string) bool {

	for i := range ssl30ciphers {
		if name == ssl30ciphers[i] {
			return true
		}
	}
	return false
}
