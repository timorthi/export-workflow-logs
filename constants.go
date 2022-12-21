package main

const (
	// GitHub Actions default env vars reference: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	envVarRepoOwner     string = "GITHUB_REPOSITORY_OWNER"
	envVarRepoFullName  string = "GITHUB_REPOSITORY"
	envVarWorkflowRunID string = "GITHUB_RUN_ID"
	envVarRunnerTempDir string = "RUNNER_TEMP"
)
