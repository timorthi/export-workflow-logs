package main

// Environment variables
// GitHub Actions default env vars reference: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
const (
	envVarGitHubServerURL string = "GITHUB_SERVER_URL"
	envVarRepoOwner       string = "GITHUB_REPOSITORY_OWNER"
	envVarRepoFullName    string = "GITHUB_REPOSITORY"
	envVarRunnerDebug     string = "RUNNER_DEBUG"
	envVarDebug           string = "DEBUG"
)

// Action inputs
const (
	inputKeyRepoToken     string = "repo-token"
	inputKeyWorkflowRunID string = "run-id"
	inputKeyDestination   string = "destination"
)

// S3 Action Inputs
const (
	inputKeyAWSAccessKeyID     string = "aws-access-key-id"
	inputKeyAWSSecretAccessKey string = "aws-secret-access-key"
	inputKeyAWSSessionToken    string = "aws-session-token"
	inputKeyAWSRegion          string = "aws-region"
	inputKeyS3BucketName       string = "s3-bucket-name"
	inputKeyS3Key              string = "s3-key"
)

// Blob Storage Action Inputs
const (
	inputKeyAzureStorageAccountName string = "azure-storage-account-name"
	inputKeyAzureStorageAccountKey  string = "azure-storage-account-key"
	inputKeyContainerName           string = "container-name"
	inputKeyBlobName                string = "blob-name"
)

// Misc constants
const (
	amazonS3Destination         string = "s3"
	azureBlobStorageDestination string = "blobstorage"
	githubDefaultBaseURL        string = "https://github.com"
)

var (
	supportedDestinations = []string{amazonS3Destination, azureBlobStorageDestination}
)
