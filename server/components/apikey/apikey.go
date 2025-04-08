package apikey

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateApiKey(length int) string {
	apikey := make([]byte, length)
	for i := range apikey {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		apikey[i] = charset[num.Int64()]
	}

	return string(apikey)
}