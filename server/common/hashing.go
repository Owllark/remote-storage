package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash256(s string) string {
	hashArr := sha256.Sum256([]byte(s))
	res := hex.EncodeToString(hashArr[:])
	return res
}
