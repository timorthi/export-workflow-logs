package main

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if os.Getenv(envVarRunnerDebug) == "1" || os.Getenv(envVarDebug) == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func main() {
	// https://go.dev/doc/go1.13#testing
	// ...testing flags are now only registered when running a test binary, and packages that call
	// flag.Parse during package initialization may cause tests to fail.
	flag.Parse()
	ctx := context.Background()

	log.Debug().Msg("Attempting to validate Action inputs")
	actionInputs, err := validateActionInputs()
	if err != nil {
		log.Fatal().Err(err).Msg("Error while validating Action inputs")
	}
	log.Info().Msg("Validated Action inputs")

	log.Debug().Int64("workflowRunID", actionInputs.workflowRunID).Msg("Attempting to get workflow run logs URL via GitHub API")
	client, err := githubClient(ctx, actionInputs.repoToken)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing GitHub client")
	}
	workflowRunLogsURL, err := getWorkflowRunLogsURLForRunID(ctx, client, actionInputs.workflowRunID)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while trying to fetch workflow run logs URL")
	}
	workflowRunLogsURLStr := workflowRunLogsURL.String()
	log.Info().Int64("workflowRunID", actionInputs.workflowRunID).Str("url", workflowRunLogsURLStr).
		Msg("Fetched URL to download workflow logs")

	log.Debug().Str("url", workflowRunLogsURLStr).Msg("Attempting to fetch workflow run logs by URL")
	workflowRunLogs, err := getResponseBodyByURL(workflowRunLogsURLStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while fetched workflow logs")
	}
	log.Info().Int64("workflowRunID", actionInputs.workflowRunID).Str("url", workflowRunLogsURLStr).
		Msg("Successfully fetched workflow run logs")

	if strings.EqualFold(actionInputs.destination, amazonS3Destination) {
		log.Debug().Msg("Attempting to upload workflow logs to S3")
		s3Client, err := s3Client(ctx, AWSConfig{
			accessKeyID:     actionInputs.s3Inputs.awsAccessKeyID,
			secretAccessKey: actionInputs.s3Inputs.awsSecretAccessKey,
			sessionToken:    actionInputs.s3Inputs.awsSessionToken,
			region:          actionInputs.s3Inputs.awsRegion,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error initializing S3 client")
		}
		err = saveToS3(ctx, s3Client, PutObjectParams{
			Bucket:   actionInputs.s3Inputs.bucketName,
			Key:      actionInputs.s3Inputs.key,
			Contents: workflowRunLogs,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error uploading workflow logs to S3")
		}
		log.Info().Str("s3BucketName", actionInputs.s3Inputs.bucketName).Str("s3Key", actionInputs.s3Inputs.key).
			Msg("Successfully saved workflow run logs to S3")
		return
	}

	if strings.EqualFold(actionInputs.destination, azureBlobStorageDestination) {
		log.Debug().Msg("Attempting to upload workflow logs to Blob Storage")
		blobStorageClient, err := blobStorageClient(AzureStorageConfig{
			storageAccountName: actionInputs.blobStorageInputs.storageAccountName,
			storageAccountKey:  actionInputs.blobStorageInputs.storageAccountKey,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error initializing Blob Storage client")
		}
		err = saveToBlobStorage(ctx, blobStorageClient, UploadBufferParams{
			ContainerName: actionInputs.blobStorageInputs.containerName,
			BlobName:      actionInputs.blobStorageInputs.blobName,
			Contents:      workflowRunLogs,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error uploading workflow logs to Blob Storage")
		}
		log.Info().Str("containerName", actionInputs.blobStorageInputs.containerName).
			Str("blobName", actionInputs.blobStorageInputs.blobName).
			Msg("Successfully saved workflow run logs to Blob Storage")
		return
	}
}
