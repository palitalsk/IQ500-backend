package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	charset    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	codeLength = 8
)

func GenerateBusinessCode() (string, error) {
	randomPart, err := generateRandomString(codeLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	// Format @XXXXXXXX
	return "@" + randomPart, nil
}

// สร้าง string แบบสุ่มจาก charset
func generateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result), nil
}
