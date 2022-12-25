package main

import (
	"io"
	"net/http"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

func downloadFileByURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpDir, err := os.MkdirTemp(getRequiredEnv(envVarGitHubWorkspace), "tmp")
	if err != nil {
		return "", err
	}

	tempFilePath := path.Join(tmpDir, tempFileName)
	log.Debug().Str("tempFilePath", tempFilePath).Msg("Creating temp file and writing contents to file")

	out, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tempFilePath, nil
}

func getRequiredEnv(envVarName string) string {
	val, exists := os.LookupEnv(envVarName)
	if !exists {
		log.Fatal().Str("envVarName", envVarName).Msg("Env var does not exist")
	}
	return val
}
