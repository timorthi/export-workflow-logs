package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

type mockPutObjectAPI func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)

func (m mockPutObjectAPI) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return m(ctx, params, optFns...)
}

func TestSaveToS3(t *testing.T) {
	ctx := context.Background()
	testContents := "hello world!"
	contentsBuf := []byte(testContents)

	testBucket := "testBucketName"
	testKey := "some/key/to/logs.zip"
	testParams := PutObjectParams{Bucket: testBucket, Key: testKey, Contents: bytes.NewBuffer(contentsBuf)}

	testCases := []struct {
		desc          string
		shouldSucceed bool
		want          string
		mockAPI       mockPutObjectAPI
	}{
		{
			desc:          "Successful PutObject call",
			shouldSucceed: true,
			want:          "",
			mockAPI: mockPutObjectAPI(
				func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
					// Assertions to test the params given to the S3 client
					assert.Equal(t, *params.Key, testKey)
					assert.Equal(t, *params.Bucket, testBucket)
					buf := new(strings.Builder)
					io.Copy(buf, params.Body)
					assert.Equal(t, buf.String(), testContents)
					return &s3.PutObjectOutput{}, nil
				}),
		},
		{
			desc:          "Failed PutObject call",
			shouldSucceed: false,
			want:          "oh no!",
			mockAPI: mockPutObjectAPI(
				func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
					return nil, errors.New("oh no!")
				}),
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := saveToS3(ctx, tC.mockAPI, testParams)

			if tC.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tC.want)
			}
		})
	}
}
