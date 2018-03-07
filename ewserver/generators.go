package ewserver

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomBytes securely generates size bytes of random data
func GenerateRandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateAPIKey for API usage
func GenerateAPIKey() (APIKey, error) {
	key, err := GenerateRandomBytes(256)
	if err != nil {
		return "", err
	}
	return APIKey(hex.EncodeToString(key)), nil
}

// GenerateRandomPassword for initial user (and allow them to change)
func GenerateRandomPassword(size int) (string, error) {
	passwd, err := GenerateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(passwd), nil
}
