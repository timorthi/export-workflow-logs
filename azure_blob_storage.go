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

// UploadBufferAPI represents the azblob SDK UploadBufferAPI call
type UploadBufferAPI interface {
	UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error)
}

// UploadBufferParams contains the required params to make an UploadBuffer call
type UploadBufferParams struct {
	ContainerName string
	BlobName      string
	Contents      *bytes.Buffer
}

// saveToBlobStorage makes an Azure Blob Storage UploadBuffer call
func saveToBlobStorage(ctx context.Context, client UploadBufferAPI, uploadParams UploadBufferParams) error {
	_, err := client.UploadBuffer(ctx, uploadParams.ContainerName, uploadParams.BlobName, uploadParams.Contents.Bytes(), nil)
	if err != nil {
		return err
	}

	return nil
}
