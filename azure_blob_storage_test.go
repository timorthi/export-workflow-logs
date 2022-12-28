package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type mockUploadBufferAPI func(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error)

func (m mockUploadBufferAPI) UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
	return m(ctx, containerName, blobName, buffer, o)
}

func TestSaveToBlobStorage(t *testing.T) {
	testContents := "hello world!"
	testParams := UploadBufferParams{
		ContainerName: "mytestcontainer",
		BlobName:      "some/blob.zip",
		Contents:      bytes.NewBuffer([]byte(testContents)),
	}

	testAPI := mockUploadBufferAPI(func(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
		t.Helper()

		if containerName != testParams.ContainerName {
			t.Fatalf("expected supplied container name to be %s, got %s", testParams.ContainerName, containerName)
		}

		if blobName != testParams.BlobName {
			t.Fatalf("expected supplied blob name to be %s, got %s", testParams.BlobName, blobName)
		}

		if string(buffer) != testContents {
			t.Fatalf("unexpected file contents: '%s'", string(buffer))
		}

		return azblob.UploadBufferResponse{}, nil
	})

	err := saveToBlobStorage(context.Background(), testAPI, testParams)
	if err != nil {
		t.Fatal(err)
	}
}
