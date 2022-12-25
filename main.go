package main

import (
	"context"
	"flag"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if os.Getenv(envVarRunnerDebug) == "1" {
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

	log.Debug().Msg("Attempting to validate Action inputs")
	err := validateActionInputs()
	if err != nil {
		log.Fatal().Err(err)
	} else {
		log.Info().Msg("Validated Action inputs")
	}

	client := githubClient()

	repoOwner := getRequiredEnv(envVarRepoOwner)
	repoName := strings.Split(getRequiredEnv(envVarRepoFullName), "/")[1]
	inputWorkflowRunID := *inputWorkflowRunIDPtr
	log.Debug().
		Str("repoName", repoName).
		Str("repoOwner", repoOwner).
		Int64("workflowRunID", inputWorkflowRunID).
		Msg("Attempting to fetch workflow run logs")

	url, _, err := client.Actions.GetWorkflowRunLogs(
		context.Background(),
		repoOwner,
		repoName,
		inputWorkflowRunID,
		true,
	)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Debug().Str("url", url.String()).Msg("Attempting to download workflow run logs by URL")
	pathToFile, err := downloadFileByURL(url.String())
	if err != nil {
		log.Fatal().Err(err)
	}
	defer os.RemoveAll(path.Dir(pathToFile))
	log.Info().Int64("workflowRunID", inputWorkflowRunID).Msg("Successfully downloaded workflow run logs")

	if strings.EqualFold(*inputDestination, "s3") {
		log.Debug().Msg("Attempting to upload workflow logs to S3")
		s3Client, err := s3Client()
		if err != nil {
			log.Fatal().Err(err)
		}
		err = saveToS3(context.Background(), s3Client, *inputS3BucketName, *inputS3Key, pathToFile)
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Info().
			Str("s3BucketName", *inputS3BucketName).
			Str("s3Key", *inputS3Key).
			Msg("Successfully saved workflow run logs to S3")
	}
}
