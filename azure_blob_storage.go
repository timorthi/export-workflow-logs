package main

import (
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
