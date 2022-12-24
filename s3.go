package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3PutObjectAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func s3Client() (*s3.Client, error) {
	os.Setenv("AWS_ACCESS_KEY_ID", *inputAWSAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *inputAWSSecretAccessKey)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(*inputAWSRegion))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return s3Client, nil
}

func saveToS3(ctx context.Context, api S3PutObjectAPI, bucket string, key string, pathToLogsFile string) error {
	logsFile, err := os.Open(pathToLogsFile)
	if err != nil {
		return err
	}
	defer logsFile.Close()

	putObjectInput := s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   logsFile,
	}
	_, err = api.PutObject(ctx, &putObjectInput)
	if err != nil {
		return err
	}

	return nil
}
