package randomip

import (
	"math/rand"
	"net"
	"time"
)

// ReservedIPv4 is a collection of reserved IPv4 addresses.
var ReservedIPv4 = []net.IPNet{
	{IP: net.IP{0, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0}},               // 0.0.0.0/8, "This" network
	{IP: net.IP{10, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0}},              // 10.0.0.0/8, CLass A private network
	{IP: net.IP{100, 64, 0, 0}, Mask: net.IPMask{255, 192, 0, 0}},          // 100.64.0.0/10, Carrier-grade NAT
	{IP: net.IP{127, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0}},             // 127.0.0.0/8, Loopback
	{IP: net.IP{169, 254, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}},         // 169.254.0.0/16, Link local
	{IP: net.IP{172, 16, 0, 0}, Mask: net.IPMask{255, 240, 0, 0}},          // 172.16.0.0/12, Class B private network
	{IP: net.IP{192, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 0}},         // 192.0.0.0/24, IETF protocol assignments
	{IP: net.IP{192, 0, 2, 0}, Mask: net.IPMask{255, 255, 255, 0}},         // 192.0.2.0/24, TEST-NET-1
	{IP: net.IP{192, 88, 99, 0}, Mask: net.IPMask{255, 255, 255, 0}},       // 192.88.99.0/24, Reserved, formerly IPv6 to IPv4
	{IP: net.IP{192, 168, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}},         // 192.168.0.0/24, Class C private network
	{IP: net.IP{198, 18, 0, 0}, Mask: net.IPMask{255, 254, 0, 0}},          // 198.18.0.0/15, Benchmarking
	{IP: net.IP{198, 51, 100, 0}, Mask: net.IPMask{255, 255, 255, 0}},      // 198.51.100.0/24, TEST-NET-2
	{IP: net.IP{203, 0, 113, 0}, Mask: net.IPMask{255, 255, 255, 0}},       // 203.0.113.0/24, TEST-NET-3
	{IP: net.IP{224, 0, 0, 0}, Mask: net.IPMask{240, 0, 0, 0}},             // 224.0.0.0/4, Multicast
	{IP: net.IP{233, 252, 0, 0}, Mask: net.IPMask{255, 255, 255, 0}},       // 233.252.0.0/24 , MCAST-TEST-NET
	{IP: net.IP{240, 0, 0, 0}, Mask: net.IPMask{240, 0, 0, 0}},             // 240.0.0.0/4, Reserved for future use
	{IP: net.IP{255, 255, 255, 255}, Mask: net.IPMask{255, 255, 255, 255}}, // 255.255.255.255/32, Broadcast
}

// ReservedIPv6 is a collection of reserved IPv6 addresses.
var ReservedIPv6 = []net.IPNet{
	{IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}}, // ::/128, Unspecified Address
	{IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}}, // ::1/128, Loopback Address
	{IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0}},     // ::ffff:0:0/96, IPv4-mapped addresses
	{IP: net.IP{0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0}},     // ::ffff:0:0:0/96, IPv4 translated addresses
	{IP: net.IP{0, 100, 255, 155, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0}},   // 64:ff9b::/96, IPv4-IPv6 Translat.
	{IP: net.IP{0, 100, 255, 155, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},               // 64:ff9b:1::/48, IPv4-IPv6 Translat.
	{IP: net.IP{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0}},                 // 100::/64, Discard-Only Address Block
	{IP: net.IP{32, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                        // 2001::/32, IETF Protocol Assignments
	{IP: net.IP{32, 1, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 240, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                       // 2001:20::/28, ORCHIDv2
	{IP: net.IP{32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                     // 2001:db8::/32, Documentation
	{IP: net.IP{32, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                            // 2002::/16, 6to4
	{IP: net.IP{252, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                             // fc00::/7, Unique-Local
	{IP: net.IP{254, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                         // fe80::/10, Link-Local Unicast
	{IP: net.IP{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.IPMask{255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},                             // ff00::/8, Multicast
}

// IsReservedIPv4 checks if the given IP address is reserved.
func IsReservedIPv4(ip net.IP) bool {
	for i := range ReservedIPv4 {
		if ReservedIPv4[i].Contains(ip) {
			return true
		}
	}
	return false
}

// IsReservedIPv6 checks if the given IP address is reserved.
func IsReservedIPv6(ip net.IP) bool {
	for i := range ReservedIPv6 {
		if ReservedIPv6[i].Contains(ip) {
			return true
		}
	}
	return false
}

// GetRandomIPv4 is return a random IPv4 address.
// The returned IP *can be* a reserved address.
func GetRandomIPv4() net.IP {

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	bytes := make([]byte, 4)

	rnd.Read(bytes)

	return net.IP{bytes[0], bytes[1], bytes[2], bytes[3]}
}

// GetRandomIPv6 is return a random IPv6 address.
// The returned IP *can be* a reserved address.
func GetRandomIPv6() net.IP {

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	bytes := make([]byte, 16)

	rnd.Read(bytes)

	return net.IP{bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7],
		bytes[8], bytes[9], bytes[10], bytes[11], bytes[12], bytes[13], bytes[14], bytes[15]}
}

// GetRandomIP is return a random IP address.
// The returned IP *can be* a reserved address.
// The version of the IP protocol is random.
func GetRandomIP() net.IP {

	n := rand.Intn(2)

	if n == 0 {
		return GetRandomIPv4()
	} else {
		return GetRandomIPv6()
	}
}

// GetPublicIPv4 is return a **non reserved** IPv4 address.
func GetPublicIPv4() net.IP {

	for {
		ip := GetRandomIPv4()

		if !IsReservedIPv4(ip) {
			return ip
		}
	}
}

// GetPublicIPv6 is return a **non reserved** IPv6 address.
func GetPublicIPv6() net.IP {

	for {
		ip := GetRandomIPv6()

		if !IsReservedIPv6(ip) {
			return ip
		}
	}
}

// GetPublicIP is return a **non reserved** IP address.
// The version of the IP protocol is random.
func GetPublicIP() net.IP {

	n := rand.Intn(2)

	if n == 0 {
		return GetPublicIPv4()
	} else {
		return GetPublicIPv6()
	}
}
