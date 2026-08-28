package attachment

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func Checksum(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func Verify(data []byte, expected string) bool {
	if strings.TrimSpace(expected) == "" {
		return false
	}
	return Checksum(data) == strings.ToLower(strings.TrimSpace(expected))
}

func ContentType(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if strings.HasSuffix(name, ".txt") {
		return "text/plain"
	}
	if strings.HasSuffix(name, ".md") {
		return "text/markdown"
	}
	if strings.HasSuffix(name, ".csv") {
		return "text/csv"
	}
	if strings.HasSuffix(name, ".json") {
		return "application/json"
	}
	return "application/octet-stream"
}
