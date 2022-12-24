package main

import (
	"context"
	"flag"
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

	tmpDir, err := os.MkdirTemp(getRequiredEnv(envVarGitHubWorkspace), "tmp")
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
	flag.Parse()
	fmt.Println(os.Args)
	client := githubClient()

	repoOwner := getRequiredEnv(envVarRepoOwner)
	repoName := strings.Split(getRequiredEnv(envVarRepoFullName), "/")[1]

	inputWorkflowRunID := *inputWorkflowRunIDPtr
	log.Printf("repoOwner:%s\nrunID:%d\nrepoName:%s\n", repoOwner, inputWorkflowRunID, repoName)

	workflow, _, err := client.Actions.GetWorkflowRunByID(context.Background(), repoOwner, repoName, inputWorkflowRunID)
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

	if strings.EqualFold(*inputDestination, "s3") {
		log.Println("Attempting to save workflow logs to S3")
		s3Client, err := s3Client()
		if err != nil {
			log.Fatal(err)
		}
		err = saveToS3(context.Background(), s3Client, *inputS3BucketName, *inputS3Key, pathToFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
