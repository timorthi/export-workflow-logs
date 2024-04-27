package main

import (
	"bytes"
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

func cloudStorageClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateObjectParams contains the required params to make an UploadBuffer call
type CreateObjectParams struct {
	BucketName string
	ObjectName string
	Contents   *bytes.Buffer
}

// saveToCloudStorage creates an object in Google Cloud Storage
func saveToCloudStorage(ctx context.Context, client *storage.Client, createObjectParams CreateObjectParams) error {
	bucket := client.Bucket(createObjectParams.BucketName)
	object := bucket.Object(createObjectParams.ObjectName)
	writer := object.NewWriter(ctx)

	_, err := fmt.Fprint(writer, createObjectParams.Contents.String())
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}
