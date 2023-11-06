package main

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3PutObjectAPI represents the AWS SDK PutObject call
type S3PutObjectAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// AWSConfig is a struct containing the AWS credentials and config needed to initialize the SDK client
type AWSConfig struct {
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
	region          string
}

func s3Client(ctx context.Context, cfg AWSConfig) (*s3.Client, error) {
	os.Setenv("AWS_ACCESS_KEY_ID", cfg.accessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.secretAccessKey)
	if len(cfg.sessionToken) > 0 {
		os.Setenv("AWS_SESSION_TOKEN", cfg.sessionToken)
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsConfig)
	return s3Client, nil
}

// PutObjectParams contains the required params to make a PutObject call
type PutObjectParams struct {
	Bucket   string
	Key      string
	Contents *bytes.Buffer
}

// saveToS3 makes an S3 PutObject call
func saveToS3(ctx context.Context, api S3PutObjectAPI, putObjectParams PutObjectParams) error {
	_, err := api.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &putObjectParams.Bucket,
		Key:    &putObjectParams.Key,
		Body:   putObjectParams.Contents,
	})
	if err != nil {
		return err
	}

	return nil
}
