package main

import (
	"context"
	"fmt"
	"os"
	"path"

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
	containerName string
	blobName      string
}

func saveToBlobStorage(ctx context.Context, client *azblob.Client, pathToLogsFile string, uploadFileParams *UploadFileParams) error {
	logsFile, err := os.Open(pathToLogsFile)
	if err != nil {
		return err
	}
	defer logsFile.Close()
	defer os.RemoveAll(path.Dir(pathToLogsFile))

	_, err = client.UploadFile(ctx, uploadFileParams.containerName, uploadFileParams.blobName, logsFile, nil)
	if err != nil {
		return err
	}

	return nil
}
