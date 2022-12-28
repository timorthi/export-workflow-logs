package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func blobStorageClient() (*azblob.Client, error) {
	storageAccountName := *inputAzureStorageAccountName
	credential, err := azblob.NewSharedKeyCredential(storageAccountName, *inputAzureStorageAccountKey)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)

	return azblob.NewClientWithSharedKeyCredential(url, credential, nil)
}

type UploadFileParams struct {
	ContainerName string
	BlobName      string
}

func saveToBlobStorage(ctx context.Context, client *azblob.Client, contents *bytes.Buffer, uploadFileParams *UploadFileParams) error {
	_, err := client.UploadBuffer(ctx, uploadFileParams.ContainerName, uploadFileParams.BlobName, contents.Bytes(), nil)
	if err != nil {
		return err
	}

	return nil
}
