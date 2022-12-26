package main

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/go-github/v48/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func githubClient() (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *inputRepoTokenPtr},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	serverURL, err := getRequiredEnv(envVarGitHubServerURL)
	if err != nil {
		return nil, err
	}
	parsedServerURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	client.BaseURL = parsedServerURL
	log.Debug().Str("serverURL", serverURL).Msgf("Set GitHub Client BaseURL to value of '%s'", envVarGitHubServerURL)

	return client, nil
}

// Uses the given workflowRunID and the GitHub Actions default environment variables to makes a GetWorkflowRunLogs call
func getWorkflowRunLogsURLForRunID(client *github.Client, workflowRunID int64) (*url.URL, error) {
	repoOwner, err := getRequiredEnv(envVarRepoOwner)
	if err != nil {
		return nil, err
	}
	repoFullName, err := getRequiredEnv(envVarRepoFullName)
	if err != nil {
		return nil, err
	}
	repoName := strings.Split(repoFullName, "/")[1]

	log.Debug().
		Str("repoName", repoName).
		Str("repoOwner", repoOwner).
		Int64("workflowRunID", workflowRunID).
		Msg("Making GetWorkflowRunLogs request")

	url, _, err := client.Actions.GetWorkflowRunLogs(
		context.Background(),
		repoOwner,
		repoName,
		workflowRunID,
		true,
	)
	if err != nil {
		return nil, err
	}

	return url, nil
}
