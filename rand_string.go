package util

import (
	"crypto/rand"
	"math/big"
)

const chars58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var charlen = big.NewInt(int64(len(chars58)))

// RandString generates a random base64 string of n characters long
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		r, _ := rand.Int(rand.Reader, charlen)
		b[i] = chars58[r.Int64()]
	}
	return string(b)
}
