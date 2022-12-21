package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func downloadFileByURL(url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpDir, err := os.MkdirTemp(GetRequiredEnv(envVarGitHubWorkspace), "tmp")
	if err != nil {
		return "", err
	}

	// Create the file
	tempFilePath := path.Join(tmpDir, tempFileName)
	log.Printf("Using path: %s", tempFilePath)
	out, err := os.Create(tempFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", err
	}

	return tempFilePath, nil
}

func main() {
	fmt.Println(os.Args)
	log.Println("Hello world!")
	client := GithubClient()

	repoOwner := GetRequiredEnv(envVarRepoOwner)
	repoName := strings.Split(GetRequiredEnv(envVarRepoFullName), "/")[1]

	log.Printf("repoOwner:%s\nrunID:%d\nrepoName:%s\n", repoOwner, runID, repoName)

	workflow, _, err := client.Actions.GetWorkflowRunByID(context.Background(), repoOwner, repoName, runID)
	log.Println(workflow)
	if err != nil {
		log.Fatal(err)
	}

	url, _, getLogsErr := client.Actions.GetWorkflowRunLogs(
		context.Background(),
		repoOwner,
		repoName,
		inputWorkflowRunID,
		true,
	)

	if getLogsErr != nil {
		log.Fatal(getLogsErr)
	}

	log.Println(url)

	pathToFile, fileDownloadErr := downloadFileByURL(url.String())
	if fileDownloadErr != nil {
		log.Fatal(fileDownloadErr)
	}
	defer os.RemoveAll(path.Dir(pathToFile))

	log.Printf("Path to file is: %s", pathToFile)
}
