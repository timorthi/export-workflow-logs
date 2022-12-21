package main

const (
	// GitHub Actions default env vars reference: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	envVarRepoOwner       string = "GITHUB_REPOSITORY_OWNER"
	envVarRepoFullName    string = "GITHUB_REPOSITORY"
	envVarGitHubWorkspace string = "GITHUB_WORKSPACE"
	envVarRepoToken       string = "INPUT_REPO-TOKEN"
	envVarWorkflowRunID   string = "INPUT_RUN-ID"
)
