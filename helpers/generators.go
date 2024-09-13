package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"
)

func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(b)

	timestamp := time.Now().UnixNano()

	apiKey := fmt.Sprintf("%s_%d", encoded, timestamp)

	return apiKey, nil
}

func GenerateOTP() (int, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	otp := n.Int64() + 100000

	return int(otp), nil
}
