package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/rs/zerolog/log"
)

func downloadFileByURL(url string) (pathToSavedFile string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	githubWorkspace, err := getRequiredEnv(envVarGitHubWorkspace)
	if err != nil {
		return "", err
	}

	tmpDir, err := os.MkdirTemp(githubWorkspace, "tmp")
	if err != nil {
		return "", err
	}

	pathToSavedFile = path.Join(tmpDir, tempFileName)
	log.Debug().Str("pathToSavedFile", pathToSavedFile).Msg("Creating temp file and writing contents to file")

	out, err := os.Create(pathToSavedFile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return pathToSavedFile, nil
}

// Returns the environment variable or an error if it is not set
func getRequiredEnv(envVarName string) (string, error) {
	val, exists := os.LookupEnv(envVarName)
	if !exists {
		return "", fmt.Errorf("env var '%s' does not exist", envVarName)

	}
	return val, nil
}
