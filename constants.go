package main

const (
	// GitHub Actions default env vars reference: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	envVarRepoOwner       string = "GITHUB_REPOSITORY_OWNER"
	envVarRepoFullName    string = "GITHUB_REPOSITORY"
	envVarGitHubWorkspace string = "GITHUB_WORKSPACE"
	envVarRunnerDebug     string = "RUNNER_DEBUG"

	inputKeyRepoToken          string = "repo-token"
	inputKeyWorkflowRunID      string = "run-id"
	inputKeyDestination        string = "destination"
	inputKeyAWSAccessKeyID     string = "aws-access-key-id"
	inputKeyAWSSecretAccessKey string = "aws-secret-access-key"
	inputKeyAWSRegion          string = "aws-region"
	inputKeyS3BucketName       string = "s3-bucket-name"
	inputKeyS3Key              string = "s3-key"

	tempFileName string = "logs.zip"
)

var (
	supportedDestinations = []string{"s3"}
)
