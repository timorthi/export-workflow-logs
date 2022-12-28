package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// Makes a GET request to the given URL and returns the response body in a buffer.
func getResponseBodyByURL(url string) (*bytes.Buffer, error) {
	log.Debug().Str("url", url).Msg("Making request to URL")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	bytesRead, err := io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("url", url).Int64("bytesRead", bytesRead).Msg("Wrote response body into buffer")

	return buf, nil
}

// Returns the environment variable or an error if it is not set
func getRequiredEnv(envVarName string) (string, error) {
	val, exists := os.LookupEnv(envVarName)
	if !exists {
		return "", fmt.Errorf("env var '%s' does not exist", envVarName)

	}
	return val, nil
}
