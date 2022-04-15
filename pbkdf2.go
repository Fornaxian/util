package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// PBKDF2Check checks the passwords of legacy pixeldrain-spring users.
func PBKDF2Check(password, legacyHash string) (valid bool, err error) {
	split := strings.Split(legacyHash, ":")

	iter, err := strconv.Atoi(split[0])
	if err != nil {
		return false, err
	}
	salt, err := hex.DecodeString(split[1])
	if err != nil {
		return false, err
	}
	hash, err := hex.DecodeString(split[2])
	if err != nil {
		return false, err
	}

	return bytes.Equal(
		hash,
		pbkdf2.Key(
			[]byte(password),
			salt,
			iter,
			len(hash),
			sha1.New,
		),
	), nil
}
