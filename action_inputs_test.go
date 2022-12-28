package main

import (
	"flag"
	"testing"
)

func TestValidateActionInputsErrorOnInvalidDestination(t *testing.T) {
	flag.Set(inputKeyDestination, "someUnsupportedDestination")
	defer flag.Set(inputKeyDestination, "")

	err := validateActionInputs()
	if err == nil {
		t.Error("Expected validateActionInputs to return error, got nil")
	}
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
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for key, value := range tC.inputValuesByKey {
				flag.Set(key, value)
				defer flag.Set(key, "") // Clean up flags by "unsetting" them after this test
			}

			err := validateActionInputs()
			if tC.shouldSucceed && err != nil {
				t.Errorf("Expected validateActionInputs to return nil, got error: %v", err)
			} else if !tC.shouldSucceed && err == nil {
				t.Error("Expected validateActionInputs to error, but it succeeded")
			}
		})
	}
}
