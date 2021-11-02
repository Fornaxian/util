package util

import (
	"net"
)

// IPMaskPrefix applies a netmask to an IPv4 or IPv6 address. Only the bits
// inside the mask are retained, the remaining bits are set to 0
func IPMaskPrefix(addr net.IP, v4mask, v6mask int) (ipMask string) {
	if addr.To4() != nil {
		// This is a IPv4 address, apply the v4 mask
		return addr.Mask(net.CIDRMask(v4mask, 32)).String()
	} else {
		// This is a IPv6 address, apply the v6 mask
		return addr.Mask(net.CIDRMask(v6mask, 128)).String()
	}
}
