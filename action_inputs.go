package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	inputRepoToken     *string = flag.String(inputKeyRepoToken, "", "GITHUB_TOKEN or a Personal Access Token")
	inputWorkflowRunID *int64  = flag.Int64(inputKeyWorkflowRunID, 0, "GitHub Actions Workflow Run ID")
	inputDestination   *string = flag.String(inputKeyDestination, "", "The service to export workflow logs to")

	// Required inputs for S3

	inputAWSAccessKeyID     *string = flag.String(inputKeyAWSAccessKeyID, "", "AWS Access Key ID")
	inputAWSSecretAccessKey *string = flag.String(inputKeyAWSSecretAccessKey, "", "AWS Secret Access Key")
	inputAWSRegion          *string = flag.String(inputKeyAWSRegion, "us-east-1", "AWS Region for the S3 bucket")
	inputS3BucketName       *string = flag.String(inputKeyS3BucketName, "", "S3 bucket name")
	inputS3Key              *string = flag.String(inputKeyS3Key, "", "S3 key")

	// Required inputs for Azure Blob Storage

	inputAzureStorageAccountName *string = flag.String(inputKeyAzureStorageAccountName, "", "Storage account name")
	inputAzureStorageAccountKey  *string = flag.String(inputKeyAzureStorageAccountKey, "", "Storage account key")
	inputContainerName           *string = flag.String(inputKeyContainerName, "", "Azure blob storage container name")
	inputBlobName                *string = flag.String(inputKeyBlobName, "", "Azure blob name")
)

type S3ActionInputs struct {
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsRegion          string
	s3BucketName       string
	s3Key              string
}

type BlobStorageActionInputs struct {
	azureStorageAccountName string
	azureStorageAccountKey  string
	containerName           string
	blobName                string
}

type ActionInputs struct {
	repoToken         string
	workflowRunID     int64
	destination       string
	s3Inputs          *S3ActionInputs
	blobStorageInputs *BlobStorageActionInputs
}

// validateActionInputs validates input combinations that cannot be checked at the action-level.
// In particular, ensures that the destination is valid and any other inputs required for that destination are present.
func validateActionInputs() (ActionInputs, error) {
	var matchedDestination string
	for _, destination := range supportedDestinations {
		if strings.EqualFold(destination, *inputDestination) {
			matchedDestination = destination
			log.Debug().Str("matchedDestination", destination).Msg("Matched input with a supported destination")
			break
		}
	}
	if matchedDestination == "" {
		return ActionInputs{}, fmt.Errorf(
			"supplied destination %s is invalid. Supported values are: %s",
			*inputDestination,
			strings.Join(supportedDestinations, ", "),
		)
	}

	var inputFlagsToAssertNotEmpty map[string]string
	var s3Inputs *S3ActionInputs
	var blobStorageInputs *BlobStorageActionInputs

	if matchedDestination == AmazonS3Destination {
		log.Debug().Msg("Validating Action inputs for S3")
		s3Inputs = &S3ActionInputs{
			awsAccessKeyID:     *inputAWSAccessKeyID,
			awsSecretAccessKey: *inputAWSSecretAccessKey,
			awsRegion:          *inputAWSRegion,
			s3BucketName:       *inputS3BucketName,
			s3Key:              *inputS3Key,
		}
		inputFlagsToAssertNotEmpty = map[string]string{
			inputKeyAWSAccessKeyID:     *inputAWSAccessKeyID,
			inputKeyAWSSecretAccessKey: *inputAWSSecretAccessKey,
			inputKeyAWSRegion:          *inputAWSRegion,
			inputKeyS3BucketName:       *inputS3BucketName,
			inputKeyS3Key:              *inputS3Key,
		}
	}

	if matchedDestination == AzureBlobStorageDestination {
		log.Debug().Msg("Validating Action inputs for Blob Storage")
		blobStorageInputs = &BlobStorageActionInputs{
			azureStorageAccountName: *inputAzureStorageAccountName,
			azureStorageAccountKey:  *inputAzureStorageAccountKey,
			containerName:           *inputContainerName,
			blobName:                *inputBlobName,
		}
		inputFlagsToAssertNotEmpty = map[string]string{
			inputKeyAzureStorageAccountName: *inputAzureStorageAccountName,
			inputKeyAzureStorageAccountKey:  *inputAzureStorageAccountKey,
			inputKeyContainerName:           *inputContainerName,
			inputKeyBlobName:                *inputBlobName,
		}
	}

	for inputName, inputValue := range inputFlagsToAssertNotEmpty {
		if len(inputValue) == 0 {
			return ActionInputs{}, fmt.Errorf("the input '%s' is required", inputName)
		}
	}

	log.Debug().Msg("Action input validation was successful")
	return ActionInputs{
		repoToken:         *inputRepoToken,
		workflowRunID:     *inputWorkflowRunID,
		destination:       matchedDestination,
		s3Inputs:          s3Inputs,
		blobStorageInputs: blobStorageInputs,
	}, nil
}
