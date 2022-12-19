package main

import (
	"context"
	"log"
	"strconv"
	"strings"
)

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
}
