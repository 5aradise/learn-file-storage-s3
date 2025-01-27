package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	errInvalidMimeType = errors.New("invalid mime type")
)

func randURLName() (string, error) {
	randSl := make([]byte, 32)
	_, err := rand.Read(randSl)
	if err != nil {
		return "", err
	}
	name := base64.RawURLEncoding.EncodeToString(randSl)
	return name, nil
}

func extractExtFromMime(mimeType string) (string, error) {
	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return "", errInvalidMimeType
	}
	return parts[1], nil
}
