package util

import (
	"net"
	"strings"
)

func IsNetError(err error) bool {
	if err == nil {
		return false
	} else if _, ok := err.(*net.OpError); ok {
		return true
	}

	return strings.HasSuffix(err.Error(), "connection reset by peer") ||
		strings.HasSuffix(err.Error(), "broken pipe") ||
		strings.HasSuffix(err.Error(), "connection timed out") ||
		strings.HasSuffix(err.Error(), "no route to host") ||
		strings.HasSuffix(err.Error(), "network is unreachable") ||
		strings.HasSuffix(err.Error(), "write: connection refused") ||
		strings.HasSuffix(err.Error(), "http2: stream closed") ||
		strings.HasSuffix(err.Error(), "client disconnected")
}
