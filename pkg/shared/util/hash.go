package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func CheckPassword(passwordInput, hashedPassword string) bool {
	hashedPasswordInput := HashPassword(passwordInput)
	return hashedPasswordInput == hashedPassword
}
