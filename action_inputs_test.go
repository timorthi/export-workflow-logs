package main

import (
	"flag"
	"testing"
)

func TestValidateActionInputsErrorOnInvalidDestination(t *testing.T) {
	flag.Set(inputKeyDestination, "someUnsupportedDestination")

	err := validateActionInputs()
	if err == nil {
		t.Error("Expected validateActionInputs to return error, got nil")
	}
}

func TestValidateActionInputsErrorOnIncompleteFlags(t *testing.T) {
	flag.Set(inputKeyRepoToken, "testRepoToken")
	flag.Set(inputKeyWorkflowRunID, "123")

	testCases := []struct {
		testDescription  string
		destination      string
		inputValuesByKey map[string]string
	}{
		{
			testDescription: "S3 destination",
			inputValuesByKey: map[string]string{
				inputKeyDestination:        "s3",
				inputKeyAWSAccessKeyID:     "abc",
				inputKeyAWSSecretAccessKey: "abc",
				inputKeyAWSRegion:          "someregion",
				// inputKeyS3BucketName intentionally excluded
				inputKeyS3Key: "some/key",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDescription, func(t *testing.T) {
			for key, value := range tc.inputValuesByKey {
				flag.Set(key, value)
			}

			err := validateActionInputs()
			if err == nil {
				t.Error("Expected validateActionInputs to return error, got nil")
			}
		})

	}
}
func TestValidateActionInputsSuccessCase(t *testing.T) {
	flag.Set(inputKeyRepoToken, "testRepoToken")
	flag.Set(inputKeyWorkflowRunID, "123")

	testCases := []struct {
		testDescription  string
		destination      string
		inputValuesByKey map[string]string
	}{
		{
			testDescription: "S3 destination",
			inputValuesByKey: map[string]string{
				inputKeyDestination:        "s3",
				inputKeyAWSAccessKeyID:     "abc",
				inputKeyAWSSecretAccessKey: "abc",
				inputKeyAWSRegion:          "someregion",
				inputKeyS3BucketName:       "my-bucket",
				inputKeyS3Key:              "some/key",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testDescription, func(t *testing.T) {
			for key, value := range tc.inputValuesByKey {
				flag.Set(key, value)
			}

			err := validateActionInputs()
			if err != nil {
				t.Errorf("Expected validateActionInputs to return nil, got error: %v", err)
			}
		})

	}
}
