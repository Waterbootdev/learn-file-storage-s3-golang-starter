package main

import (
	"crypto/rand"
	"encoding/base64"
)

func randomId(n int) (string, error) {
	randomBytes := make([]byte, n)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes), nil
}
