package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// GenerateRandomString generates a random string of the specified length.
func GenerateRandomString(length int) string {
	byteLen := length * 6 / 8
	bytes := make([]byte, byteLen)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// CreateTruncatedString creates a truncated string of the specified length.
func CreateTruncatedString(str string, length int) string {
	strLen := len(str)
	if strLen > length {
		strLen = length
	}
	if strLen == 0 {
		return ""
	}
	return str[:strLen] + "..."
}

// EncodeBase64 encodes a string into a Base64 string.
func EncodeBase64(s string) string {
	b := []byte(s)
	return base64.URLEncoding.EncodeToString(b)
}

// DecodeBase64 decodes a Base64 string into a string.
func DecodeBase64(s string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}
