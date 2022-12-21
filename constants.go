package main

import "flag"

const (
	// GitHub Actions default env vars reference: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	envVarRepoOwner       string = "GITHUB_REPOSITORY_OWNER"
	envVarRepoFullName    string = "GITHUB_REPOSITORY"
	envVarGitHubWorkspace string = "GITHUB_WORKSPACE"

	tempFileName string = "logs.zip"
)

var (
	inputRepoToken     string = *flag.String("repo-token", "", "GITHUB_TOKEN or a Personal Access Token")
	inputWorkflowRunID int64  = *flag.Int64("run-id", 0, "GitHub Actions Workflow Run ID")
)
