package main

import (
	"context"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func blobStorageClient() (*azblob.Client, error) {
	credential, err := azblob.NewSharedKeyCredential("", "")
	if err != nil {
		return nil, err
	}

	url := "https://<StorageAccountName>.blob.core.windows.net/" //replace <StorageAccountName> with your Azure storage account name

	return azblob.NewClientWithSharedKeyCredential(url, credential, nil)
}

func saveToBlobStorage(ctx context.Context, client *azblob.Client, pathToLogsFile string, containerName string, blobName string) error {
	logsFile, err := os.Open(pathToLogsFile)
	if err != nil {
		return err
	}
	defer logsFile.Close()
	defer os.RemoveAll(path.Dir(pathToLogsFile))

	_, err = client.UploadFile(ctx, containerName, blobName, logsFile, nil)
	if err != nil {
		return err
	}

	return nil
}
