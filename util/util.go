package util

import "encoding/base64"

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
