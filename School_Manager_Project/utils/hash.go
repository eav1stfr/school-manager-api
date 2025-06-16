package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

func Hash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", ErrorGeneratingSaltForHashing
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}

func VerifyPassword(passwordFromDB, providedPassword string) (bool, error) {
	parts := strings.Split(passwordFromDB, ".")
	if len(parts) != 2 {
		return false, InvalidEncodedHashFormat
	}
	saltBase64 := parts[0]
	hashedPasswordBase64 := parts[1]
	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false, FailedToDecodeSalt
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(hashedPasswordBase64)
	if err != nil {
		return false, FailedToDecodeHashError
	}
	hash := argon2.IDKey([]byte(providedPassword), salt, 1, 64*1024, 4, 32)

	if len(hash) != len(hashedPassword) || subtle.ConstantTimeCompare(hash, hashedPassword) != 1 {
		return false, IncorrectPasswordError
	}
	return true, nil
}
