package ewserver

import (
	"crypto/rand"
	"encoding/base64"
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

// GenerateRandomString generates a random string
func GenerateRandomString(size int) (string, error) {
	b, err := GenerateRandomBytes(size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// GenerateAPIKey for API usage
func GenerateAPIKey() (APIKey, error) {
	key, err := GenerateRandomString(64)
	if err != nil {
		return "", err
	}
	return APIKey(key), nil
}
