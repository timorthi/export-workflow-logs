package main

import (
	"log"
	"os"
)

func GetRequiredEnv(envVarName string) string {
	val, exists := os.LookupEnv(envVarName)
	if !exists {
		log.Fatalf("GetRequiredEnv: Env var '%s' does not exist", envVarName)
	}
	return val
}
