package util

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

// RemoteAddress resolves the address of a HTTP request to the real remote
// address. It takes in account proxies if the request is originated from a
// private IP range
func RemoteAddress(r *http.Request) (addr string) {
	var err error
	addr = r.RemoteAddr

	if addr == "@" { // Support unix domain sockets
		addr = "127.0.0.1:80"
	}

	addr, _, err = net.SplitHostPort(addr)
	if err != nil {
		panic(fmt.Errorf("failed to parse request RemoteAddr: %s", err))
	}

	if AddressIsLocal(addr) {
		// This IP is in a private range, so we can trust its proxy headers,
		// check if it has any
		if h := r.Header.Get("X-Real-IP"); h != "" {
			return h
		} else if h := r.Header.Get("X-Forwarded-For"); h != "" {
			return strings.Split(h, ", ")[0]
		}
	}

	return addr
}

// AddressIsLocal checks if an IP address falls on a private subnet
func AddressIsLocal(ip string) bool {
	if ip == "@" { // Support unix domain sockets
		ip = "127.0.0.1"
	}

	ipp := net.ParseIP(ip)
	if ipp == nil {
		return false
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ipp) {
			return true
		}
	}
	return false
}
