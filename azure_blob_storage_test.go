package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/stretchr/testify/assert"
)

type mockUploadBufferAPI func(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error)

func (m mockUploadBufferAPI) UploadBuffer(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
	return m(ctx, containerName, blobName, buffer, o)
}

func TestSaveToBlobStorage(t *testing.T) {
	ctx := context.Background()
	testContents := "hello world!"
	testParams := UploadBufferParams{
		ContainerName: "mytestcontainer",
		BlobName:      "some/blob.zip",
		Contents:      bytes.NewBuffer([]byte(testContents)),
	}

	testCases := []struct {
		desc          string
		shouldSucceed bool
		want          string
		mockAPI       mockUploadBufferAPI
	}{
		{
			desc:          "Successful UploadBuffer call",
			shouldSucceed: true,
			want:          "",
			mockAPI: mockUploadBufferAPI(
				func(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
					assert.Equal(t, containerName, testParams.ContainerName)
					assert.Equal(t, blobName, testParams.BlobName)
					assert.Equal(t, string(buffer), testContents)
					return azblob.UploadBufferResponse{}, nil
				}),
		},
		{
			desc:          "Failed UploadBuffer call",
			shouldSucceed: false,
			want:          "some error",
			mockAPI: mockUploadBufferAPI(
				func(ctx context.Context, containerName string, blobName string, buffer []byte, o *azblob.UploadBufferOptions) (azblob.UploadBufferResponse, error) {
					return azblob.UploadBufferResponse{}, errors.New("some error")
				}),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := saveToBlobStorage(ctx, tC.mockAPI, testParams)

			if tC.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tC.want)
			}
		})
	}
}
