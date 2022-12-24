package main

import (
	"os"

	"github.com/rs/zerolog/log"
)

func getRequiredEnv(envVarName string) string {
	val, exists := os.LookupEnv(envVarName)
	if !exists {
		log.Fatal().Str("envVarName", envVarName).Msg("Env var does not exist")
	}
	return val
}
