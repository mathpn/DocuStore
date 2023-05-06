package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func hashDocument(text string) string {
	hash := sha256.Sum256([]byte(text))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}
