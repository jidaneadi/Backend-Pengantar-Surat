package utils

import (
	"Backend_TA/config"
	"crypto/sha256"
	"encoding/hex"
)

func EncryptHash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))

	key := config.RenderEnv("KEY_HASH")
	hasher.Write([]byte(key))

	hash := hasher.Sum(nil)

	hashString := hex.EncodeToString(hash)

	return hashString
}

// func EncryptAES(data string) string {

// }
