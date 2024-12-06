package util

import (
	"errors"
	"net"
	"strings"
	"syscall"
)

func IsNetError(err error) bool {
	if err == nil {
		return false
	} else if _, ok := err.(*net.OpError); ok {
		return true
	}

	return errors.Is(err, syscall.ECONNRESET) || // connection reset by peer
		errors.Is(err, syscall.EPIPE) || // broken pipe
		errors.Is(err, syscall.EHOSTUNREACH) || // no route to host
		errors.Is(err, syscall.ENETUNREACH) || // network is unreachable
		errors.Is(err, syscall.ENETDOWN) || // network is down
		errors.Is(err, syscall.ECONNREFUSED) || // connection refused
		errors.Is(err, syscall.ETIMEDOUT) || // connection timed out
		strings.HasSuffix(err.Error(), "; CANCEL") || // http2.ErrCodeCancel
		strings.HasSuffix(err.Error(), "; PROTOCOL_ERROR") || // http2.ErrCodeProtocol
		strings.HasSuffix(err.Error(), "http2: stream closed") || // http2.ErrCodeStreamClosed
		strings.HasSuffix(err.Error(), "http2: request body closed due to handler exiting") || // http2.errHandlerComplete
		strings.HasSuffix(err.Error(), "client disconnected") || // http.http2errClientDisconnected
		strings.HasSuffix(err.Error(), "H3_REQUEST_CANCELLED")
}
