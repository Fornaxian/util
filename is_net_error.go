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
		strings.HasSuffix(err.Error(), "no route to host")
}
