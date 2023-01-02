package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// Required Action inputs
var (
	inputRepoToken     *string = flag.String(inputKeyRepoToken, "", "GITHUB_TOKEN or a Personal Access Token")
	inputWorkflowRunID *int64  = flag.Int64(inputKeyWorkflowRunID, 0, "GitHub Actions Workflow Run ID")
	inputDestination   *string = flag.String(inputKeyDestination, "", "The service to export workflow logs to")
)

// Required inputs for S3
var (
	inputAWSAccessKeyID     *string = flag.String(inputKeyAWSAccessKeyID, "", "AWS Access Key ID")
	inputAWSSecretAccessKey *string = flag.String(inputKeyAWSSecretAccessKey, "", "AWS Secret Access Key")
	inputAWSRegion          *string = flag.String(inputKeyAWSRegion, "us-east-1", "AWS Region for the S3 bucket")
	inputS3BucketName       *string = flag.String(inputKeyS3BucketName, "", "S3 bucket name")
	inputS3Key              *string = flag.String(inputKeyS3Key, "", "S3 key")
)

// Required inputs for Azure Blob Storage
var (
	inputAzureStorageAccountName *string = flag.String(inputKeyAzureStorageAccountName, "", "Storage account name")
	inputAzureStorageAccountKey  *string = flag.String(inputKeyAzureStorageAccountKey, "", "Storage account key")
	inputContainerName           *string = flag.String(inputKeyContainerName, "", "Azure blob storage container name")
	inputBlobName                *string = flag.String(inputKeyBlobName, "", "Azure blob name")
)

// S3ActionInputs contains inputs required for the `s3` destination
type S3ActionInputs struct {
	awsAccessKeyID     string
	awsSecretAccessKey string
	awsRegion          string
	bucketName         string
	key                string
}

// BlobStorageActionInputs contains inputs required for the `blobstorage` destination
type BlobStorageActionInputs struct {
	storageAccountName string
	storageAccountKey  string
	containerName      string
	blobName           string
}

// ActionInputs contains all the pertinent inputs for this GitHub Action. For a given destination, its corresponding
// struct field (e.g. s3Inputs for the `s3` destination) is assumed to be non-nil.
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
			bucketName:         *inputS3BucketName,
			key:                *inputS3Key,
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
			storageAccountName: *inputAzureStorageAccountName,
			storageAccountKey:  *inputAzureStorageAccountKey,
			containerName:      *inputContainerName,
			blobName:           *inputBlobName,
		}
		inputFlagsToAssertNotEmpty = map[string]string{
			inputKeyAzureStorageAccountName: *inputAzureStorageAccountName,
			inputKeyAzureStorageAccountKey:  *inputAzureStorageAccountKey,
			inputKeyContainerName:           *inputContainerName,
			inputKeyBlobName:                *inputBlobName,
		}
	}

	var emptyInputs []string
	for inputName, inputValue := range inputFlagsToAssertNotEmpty {
		if len(inputValue) == 0 {
			emptyInputs = append(emptyInputs, inputName)
		}
	}
	if len(emptyInputs) > 0 {
		sort.Strings(emptyInputs)
		return ActionInputs{}, fmt.Errorf("the following inputs are required: %s", strings.Join(emptyInputs, ", "))
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
