package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256hashing(s string) string {
	hashArr := sha256.Sum256([]byte(s))
	res := hex.EncodeToString(hashArr[:])
	return res
}
