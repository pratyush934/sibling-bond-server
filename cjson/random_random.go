package cjson

import (
	"crypto/rand"
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
