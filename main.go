package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func downloadFileByURL(url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	tempFilePath := fmt.Sprintf("%s/tmp/logs.zip", GetRequiredEnv(envVarGitHubWorkspace))
	log.Printf("Using path: %s", tempFilePath)
	out, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	log.Println("Hello world!")
	client := GithubClient()

	runIDString := GetRequiredEnv(envVarWorkflowRunID)
	runID, parseErr := strconv.ParseInt(runIDString, 10, 64)
	if parseErr != nil {
		log.Fatalf("Could not convert runID '%s' to int64", runIDString)
	}

	url, response, getLogsErr := client.Actions.GetWorkflowRunLogs(
		context.Background(),
		GetRequiredEnv(envVarRepoOwner),
		strings.Split(GetRequiredEnv(envVarRepoFullName), "/")[1],
		runID,
		true,
	)

	if getLogsErr != nil {
		log.Fatal(getLogsErr)
	}

	log.Println(url)
	log.Println(response)

	fileDownloadErr := downloadFileByURL(url.String())
	if fileDownloadErr != nil {
		log.Fatal(fileDownloadErr)
	}

	files, err := os.ReadDir(fmt.Sprintf("%s/tmp", GetRequiredEnv(envVarGitHubWorkspace)))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}
