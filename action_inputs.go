package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	inputRepoTokenPtr     *string = flag.String(inputKeyRepoToken, "", "GITHUB_TOKEN or a Personal Access Token")
	inputWorkflowRunIDPtr *int64  = flag.Int64(inputKeyWorkflowRunID, 0, "GitHub Actions Workflow Run ID")
	inputDestination      *string = flag.String(inputKeyDestination, "", "The service to export workflow logs to")

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

// Validates input combinations that cannot be checked at the action-level.
// In particular, ensures that the destination is valid and any other inputs
// required for that destination are present.
func validateActionInputs() error {
	var matchedDestination string
	for _, destination := range supportedDestinations {
		if strings.EqualFold(destination, *inputDestination) {
			matchedDestination = destination
			log.Debug().Str("matchedDestination", destination).Msg("Matched input with a supported destination")
			break
		}
	}
	if matchedDestination == "" {
		return fmt.Errorf(
			"supplied destination %s is invalid. Supported values are: %s",
			*inputDestination,
			strings.Join(supportedDestinations, ", "),
		)
	}

	var inputFlagsToAssertNotEmpty map[string]string
	if matchedDestination == "s3" {
		log.Debug().Msg("Validating Action inputs for S3")
		inputFlagsToAssertNotEmpty = map[string]string{
			inputKeyAWSAccessKeyID:     *inputAWSAccessKeyID,
			inputKeyAWSSecretAccessKey: *inputAWSSecretAccessKey,
			inputKeyAWSRegion:          *inputAWSRegion,
			inputKeyS3BucketName:       *inputS3BucketName,
			inputKeyS3Key:              *inputS3Key,
		}
	}

	if matchedDestination == "blobstorage" {
		log.Debug().Msg("Validating Action inputs for Blob Storage")
		inputFlagsToAssertNotEmpty = map[string]string{
			inputKeyAzureStorageAccountName: *inputAzureStorageAccountName,
			inputKeyAzureStorageAccountKey:  *inputAzureStorageAccountKey,
			inputKeyContainerName:           *inputContainerName,
			inputKeyBlobName:                *inputBlobName,
		}
	}

	for inputName, inputValue := range inputFlagsToAssertNotEmpty {
		if len(inputValue) == 0 {
			return fmt.Errorf("the input '%s' is required", inputName)
		}
	}

	log.Debug().Msg("Action input validation was successful")
	return nil
}
