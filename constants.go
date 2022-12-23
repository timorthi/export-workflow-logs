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
	inputRepoTokenPtr     *string = flag.String("repo-token", "", "GITHUB_TOKEN or a Personal Access Token")
	inputWorkflowRunIDPtr *int64  = flag.Int64("run-id", 0, "GitHub Actions Workflow Run ID")
	inputDestination      *string = flag.String("destination", "", "The service to export workflow logs to")

	inputAWSAccessKeyID     *string = flag.String("aws-access-key-id", "", "AWS Access Key ID")
	inputAWSSecretAccessKey *string = flag.String("aws-secret-access-key", "", "AWS Secret Access Key")
	inputS3BucketName       *string = flag.String("s3-bucket-name", "", "S3 bucket name")
	inputS3Key              *string = flag.String("s3-key", "", "S3 key")
)
