package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateActionInputsErrorOnInvalidDestination(t *testing.T) {
	flag.Set(inputKeyDestination, "someUnsupportedDestination")
	defer flag.Set(inputKeyDestination, "")

	err := validateActionInputs()
	assert.ErrorContains(t, err, "supplied destination someUnsupportedDestination is invalid")
}

func TestValidateActionInputs(t *testing.T) {
	flag.Set(inputKeyRepoToken, "testRepoToken")
	flag.Set(inputKeyWorkflowRunID, "123")
	defer flag.Set(inputKeyRepoToken, "")
	defer flag.Set(inputKeyWorkflowRunID, "")

	testCases := []struct {
		desc             string
		shouldSucceed    bool
		inputValuesByKey map[string]string
		want             string
	}{
		{
			desc:          "S3 destination success case",
			shouldSucceed: true,
			inputValuesByKey: map[string]string{
				inputKeyDestination:        "s3",
				inputKeyAWSAccessKeyID:     "abc",
				inputKeyAWSSecretAccessKey: "abc",
				inputKeyAWSRegion:          "someregion",
				inputKeyS3BucketName:       "my-bucket",
				inputKeyS3Key:              "some/key",
			},
			want: "",
		},
		{
			desc:          "S3 destination failure case",
			shouldSucceed: false,
			inputValuesByKey: map[string]string{
				inputKeyDestination:        "s3",
				inputKeyAWSAccessKeyID:     "abc",
				inputKeyAWSSecretAccessKey: "abc",
				inputKeyAWSRegion:          "someregion",
				// inputKeyS3BucketName intentionally excluded
				inputKeyS3Key: "some/key",
			},
			want: inputKeyS3BucketName,
		},
		{
			desc:          "Blob Storage destination success case",
			shouldSucceed: true,
			inputValuesByKey: map[string]string{
				inputKeyDestination:             "blobstorage",
				inputKeyAzureStorageAccountName: "mystorageaccount",
				inputKeyAzureStorageAccountKey:  "myaccesskey",
				inputKeyContainerName:           "my-container",
				inputKeyBlobName:                "logs.zip",
			},
			want: "",
		},
		{
			desc:          "Blob Storage destination failure case",
			shouldSucceed: false,
			inputValuesByKey: map[string]string{
				inputKeyDestination:             "blobstorage",
				inputKeyAzureStorageAccountName: "mystorageaccount",
				// inputKeyAzureStorageAccountKey intentionally excluded
				inputKeyContainerName: "my-container",
				inputKeyBlobName:      "logs.zip",
			},
			want: inputKeyAzureStorageAccountKey,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for key, value := range tC.inputValuesByKey {
				flag.Set(key, value)
				defer flag.Set(key, "") // Clean up flags by "unsetting" them after this test
			}

			err := validateActionInputs()
			if tC.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, fmt.Sprintf("the input '%s' is required", tC.want))
			}
		})
	}
}
