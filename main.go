package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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
	tempFilePath := fmt.Sprintf("%s/logs.zip", GetRequiredEnv(envVarRunnerTempDir))
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

	url, response, err := client.Actions.GetWorkflowRunLogs(
		context.Background(),
		GetRequiredEnv(envVarRepoOwner),
		strings.Split(GetRequiredEnv(envVarRepoFullName), "/")[1],
		runID,
		true,
	)

	log.Println(url)
	log.Println(response)
	log.Println(err)

	downloadFileByURL(url.String())

	files, err := ioutil.ReadDir(GetRequiredEnv(envVarRunnerTempDir))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}
