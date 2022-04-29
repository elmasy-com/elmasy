package ciphersuite

var CipherSuites = []CipherSuite{

	{[]byte{0x00, 0x1C}, "SSL_FORTEZZA_KEA_WITH_NULL_SHA", []uint16{SSL30}},
	{[]byte{0x00, 0x1D}, "SSL_FORTEZZA_KEA_WITH_FORTEZZA_CBC_SHA", []uint16{SSL30}},

	{[]byte{0x00, 0x00}, "TLS_NULL_WITH_NULL_NULL", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x01}, "TLS_RSA_WITH_NULL_MD5", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x02}, "TLS_RSA_WITH_NULL_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x03}, "TLS_RSA_EXPORT_WITH_RC4_40_MD5", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x04}, "TLS_RSA_WITH_RC4_128_MD5", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x05}, "TLS_RSA_WITH_RC4_128_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x06}, "TLS_RSA_EXPORT_WITH_RC2_CBC_40_MD5", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x07}, "TLS_RSA_WITH_IDEA_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x08}, "TLS_RSA_EXPORT_WITH_DES40_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x09}, "TLS_RSA_WITH_DES_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x0A}, "TLS_RSA_WITH_3DES_EDE_CBC_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x0B}, "TLS_DH_DSS_EXPORT_WITH_DES40_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x0C}, "TLS_DH_DSS_WITH_DES_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x0D}, "TLS_DH_DSS_WITH_3DES_EDE_CBC_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x0E}, "TLS_DH_RSA_EXPORT_WITH_DES40_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x0F}, "TLS_DH_RSA_WITH_DES_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x10}, "TLS_DH_RSA_WITH_3DES_EDE_CBC_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x11}, "TLS_DHE_DSS_EXPORT_WITH_DES40_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x12}, "TLS_DHE_DSS_WITH_DES_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x13}, "TLS_DHE_DSS_WITH_3DES_EDE_CBC_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x14}, "TLS_DHE_RSA_EXPORT_WITH_DES40_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x15}, "TLS_DHE_RSA_WITH_DES_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x16}, "TLS_DHE_RSA_WITH_3DES_EDE_CBC_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x17}, "TLS_DH_anon_EXPORT_WITH_RC4_40_MD5", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x18}, "TLS_DH_anon_WITH_RC4_128_MD5", []uint16{SSL30, TLS10, TLS11, TLS12}},
	{[]byte{0x00, 0x19}, "TLS_DH_anon_EXPORT_WITH_DES40_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x1A}, "TLS_DH_anon_WITH_DES_CBC_SHA", []uint16{SSL30, TLS10, TLS11}},
	{[]byte{0x00, 0x1B}, "TLS_DH_anon_WITH_3DES_EDE_CBC_SHA", []uint16{SSL30, TLS10, TLS11, TLS12}},

	// TLSKRB
	{[]byte{0x00, 0x1E}, "TLS_KRB5_WITH_DES_CBC_SHA", []uint16{SSL30, TLS11}},
	{[]byte{0x00, 0x1F}, "TLS_KRB5_WITH_3DES_EDE_CBC_SHA", []uint16{TLS11}},
	{[]byte{0x00, 0x20}, "TLS_KRB5_WITH_RC4_128_SHA", []uint16{TLS11}},
	{[]byte{0x00, 0x21}, "TLS_KRB5_WITH_IDEA_CBC_SHA", []uint16{TLS11}},
	{[]byte{0x00, 0x22}, "TLS_KRB5_WITH_DES_CBC_MD5", []uint16{TLS11}},
	{[]byte{0x00, 0x23}, "TLS_KRB5_WITH_3DES_EDE_CBC_MD5", []uint16{TLS11}},
	{[]byte{0x00, 0x24}, "TLS_KRB5_WITH_RC4_128_MD5", []uint16{TLS11}},
	{[]byte{0x00, 0x25}, "TLS_KRB5_WITH_IDEA_CBC_MD5", []uint16{TLS11}},

	// TLSKRB, MUST NOT negotiate
	{[]byte{0x00, 0x26}, "TLS_KRB5_EXPORT_WITH_DES_CBC_40_SHA", []uint16{TLS11}},
	{[]byte{0x00, 0x27}, "TLS_KRB5_EXPORT_WITH_RC2_CBC_40_SHA", []uint16{TLS11}},
	{[]byte{0x00, 0x28}, "TLS_KRB5_EXPORT_WITH_RC4_40_SHA", []uint16{TLS11}},
	{[]byte{0x00, 0x29}, "TLS_KRB5_EXPORT_WITH_DES_CBC_40_MD5", []uint16{TLS11}},
	{[]byte{0x00, 0x2A}, "TLS_KRB5_EXPORT_WITH_RC2_CBC_40_MD5", []uint16{TLS11}},
	{[]byte{0x00, 0x2B}, "TLS_KRB5_EXPORT_WITH_RC4_40_MD5", []uint16{TLS11}},

	// TLSAES
	{[]byte{0x00, 0x2F}, "TLS_RSA_WITH_AES_128_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x30}, "TLS_DH_DSS_WITH_AES_128_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x31}, "TLS_DH_RSA_WITH_AES_128_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x32}, "TLS_DHE_DSS_WITH_AES_128_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x33}, "TLS_DHE_RSA_WITH_AES_128_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x34}, "TLS_DH_anon_WITH_AES_128_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x35}, "TLS_RSA_WITH_AES_256_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x36}, "TLS_DH_DSS_WITH_AES_256_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x37}, "TLS_DH_RSA_WITH_AES_256_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x38}, "TLS_DHE_DSS_WITH_AES_256_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x39}, "TLS_DHE_RSA_WITH_AES_256_CBC_SHA", []uint16{TLS11, TLS12}},
	{[]byte{0x00, 0x3A}, "TLS_DH_anon_WITH_AES_256_CBC_SHA", []uint16{TLS11, TLS12}},

	{[]byte{0x00, 0x3B}, "TLS_RSA_WITH_NULL_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x3C}, "TLS_RSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x3D}, "TLS_RSA_WITH_AES_256_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x3E}, "TLS_DH_DSS_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x3F}, "TLS_DH_RSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x40}, "TLS_DHE_DSS_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x67}, "TLS_DHE_RSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x68}, "TLS_DH_DSS_WITH_AES_256_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x69}, "TLS_DH_RSA_WITH_AES_256_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x6A}, "TLS_DHE_DSS_WITH_AES_256_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x6B}, "TLS_DHE_RSA_WITH_AES_256_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x6C}, "TLS_DH_anon_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x6D}, "TLS_DH_anon_WITH_AES_256_CBC_SHA256", []uint16{TLS12}},

	{[]byte{0x00, 0x9C}, "TLS_RSA_WITH_AES_128_GCM_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x9D}, "TLS_RSA_WITH_AES_256_GCM_SHA384", []uint16{TLS12}},
	{[]byte{0x00, 0x9E}, "TLS_DHE_RSA_WITH_AES_128_GCM_SHA256", []uint16{TLS12}},
	{[]byte{0x00, 0x9F}, "TLS_DHE_RSA_WITH_AES_256_GCM_SHA384", []uint16{TLS12}},

	{[]byte{0xC0, 0x09}, "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA", []uint16{TLS12}},
	{[]byte{0xC0, 0x0A}, "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA", []uint16{TLS12}},
	{[]byte{0xC0, 0x13}, "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", []uint16{TLS12}},
	{[]byte{0xC0, 0x14}, "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA", []uint16{TLS12}},

	// RFC 5289
	{[]byte{0xC0, 0x23}, "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x24}, "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x25}, "TLS_ECDH_ECDSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x26}, "TLS_ECDH_ECDSA_WITH_AES_256_CBC_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x27}, "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x28}, "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x29}, "TLS_ECDH_RSA_WITH_AES_128_CBC_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x2A}, "TLS_ECDH_RSA_WITH_AES_256_CBC_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x2B}, "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x2C}, "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x2D}, "TLS_ECDH_ECDSA_WITH_AES_128_GCM_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x2E}, "TLS_ECDH_ECDSA_WITH_AES_256_GCM_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x2F}, "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x30}, "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", []uint16{TLS12}},
	{[]byte{0xC0, 0x31}, "TLS_ECDH_RSA_WITH_AES_128_GCM_SHA256", []uint16{TLS12}},
	{[]byte{0xC0, 0x32}, "TLS_ECDH_RSA_WITH_AES_256_GCM_SHA384", []uint16{TLS12}},

	// RFC 7905
	{[]byte{0xCC, 0xA8}, "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
	{[]byte{0xCC, 0xA9}, "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
	{[]byte{0xCC, 0xAA}, "TLS_DHE_RSA_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
	{[]byte{0xCC, 0xAB}, "TLS_PSK_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
	{[]byte{0xCC, 0xAC}, "TLS_ECDHE_PSK_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
	{[]byte{0xCC, 0xAD}, "TLS_DHE_PSK_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
	{[]byte{0xCC, 0xAE}, "TLS_RSA_PSK_WITH_CHACHA20_POLY1305_SHA256", []uint16{TLS12}},
}