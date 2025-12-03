package utils

import (
	"crypto/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateAPIKey(blockLength int, blocks int) string {
	var sb strings.Builder

	for i := 0; i < blocks; i++ {
		if i > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(randomString(blockLength))
	}

	return sb.String()
}

func randomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	var sb strings.Builder
	for _, v := range b {
		sb.WriteByte(charset[int(v)%len(charset)])
	}

	return sb.String()
}
