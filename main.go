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
	err := validateActionInputs()
	if err != nil {
		log.Fatal().Err(err).Msg("Error while validating Action inputs")
	}
	log.Info().Msg("Validated Action inputs")

	workflowRunID := *inputWorkflowRunIDPtr
	log.Debug().Int64("workflowRunID", workflowRunID).Msg("Attempting to get workflow run logs URL via GitHub API")
	client, err := githubClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing GitHub client")
	}
	workflowRunLogsURL, err := getWorkflowRunLogsURLForRunID(ctx, client, workflowRunID)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while trying to fetch workflow run logs URL")
	}
	workflowRunLogsURLStr := workflowRunLogsURL.String()
	log.Info().
		Int64("workflowRunID", workflowRunID).
		Str("url", workflowRunLogsURLStr).
		Msg("Fetched URL to download workflow logs")

	log.Debug().Str("url", workflowRunLogsURLStr).Msg("Attempting to download workflow run logs by URL")
	workflowRunLogs, err := getResponseBodyByURL(workflowRunLogsURLStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while downloading workflow logs")
	}
	log.Info().
		Int64("workflowRunID", workflowRunID).
		Str("url", workflowRunLogsURLStr).
		Msg("Successfully downloaded workflow run logs")

	if strings.EqualFold(*inputDestination, "s3") {
		log.Debug().Msg("Attempting to upload workflow logs to S3")
		s3Client, err := s3Client()
		if err != nil {
			log.Fatal().Err(err).Msg("Error initializing S3 client")
		}
		params := PutObjectParams{
			Bucket: *inputS3BucketName,
			Key:    *inputS3Key,
		}
		err = saveToS3(ctx, s3Client, workflowRunLogs, params)
		if err != nil {
			log.Fatal().Err(err).Msg("Error uploading workflow logs to S3")
		}
		log.Info().
			Str("s3BucketName", *inputS3BucketName).
			Str("s3Key", *inputS3Key).
			Msg("Successfully saved workflow run logs to S3")
	} else if strings.EqualFold(*inputDestination, "blobstorage") {
		log.Debug().Msg("Attempting to upload workflow logs to Blob Storage")
		blobStorageClient, err := blobStorageClient()
		if err != nil {
			log.Panic().Err(err).Msg("Error initializing Blob Storage client")
		}
		err = saveToBlobStorage(ctx, blobStorageClient, UploadParams{
			ContainerName: *inputContainerName,
			BlobName:      *inputBlobName,
			Contents:      workflowRunLogs,
		})
		if err != nil {
			log.Panic().Err(err).Msg("Error uploading workflow logs to Blob Storage")
		}
		log.Info().
			Str("containerName", *inputContainerName).
			Str("blobName", *inputBlobName).
			Msg("Successfully saved workflow run logs to Blob Storage")
	}
}
