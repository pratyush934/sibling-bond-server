package cjson

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/rs/zerolog/log"
)

func CreateRandomToken() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Err(err).Msg("Issue in CreateRandomToken")
		return ""
	}
	return hex.EncodeToString(bytes)
}

func HashRawToken(raw string) string {
	sum224 := sha256.Sum224([]byte(raw))
	return hex.EncodeToString(sum224[:])
}
