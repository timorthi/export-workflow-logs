package main

const (
	// GitHub Actions default env vars reference: https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	envVarGitHubServerURL string = "GITHUB_SERVER_URL"
	envVarRepoOwner       string = "GITHUB_REPOSITORY_OWNER"
	envVarRepoFullName    string = "GITHUB_REPOSITORY"
	envVarRunnerDebug     string = "RUNNER_DEBUG"
	envVarDebug           string = "DEBUG"

	tempFileName         string = "logs.zip"
	githubDefaultBaseURL string = "https://github.com"

	inputKeyRepoToken     string = "repo-token"
	inputKeyWorkflowRunID string = "run-id"
	inputKeyDestination   string = "destination"

	inputKeyAWSAccessKeyID     string = "aws-access-key-id"
	inputKeyAWSSecretAccessKey string = "aws-secret-access-key"
	inputKeyAWSRegion          string = "aws-region"
	inputKeyS3BucketName       string = "s3-bucket-name"
	inputKeyS3Key              string = "s3-key"

	inputKeyAzureStorageAccountName string = "azure-storage-account-name"
	inputKeyAzureStorageAccountKey  string = "azure-storage-account-key"
	inputKeyContainerName           string = "container-name"
	inputKeyBlobName                string = "blob-name"
)

const (
	AmazonS3Destination         string = "s3"
	AzureBlobStorageDestination string = "blobstorage"
)

var (
	supportedDestinations = []string{AmazonS3Destination, AzureBlobStorageDestination}
)
