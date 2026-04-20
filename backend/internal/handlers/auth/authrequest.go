package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateState() string {
	bytes := make([]byte, 36)
	rand.Read(bytes)

	return base64.RawURLEncoding.EncodeToString(bytes)
}
