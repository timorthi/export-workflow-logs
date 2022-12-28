package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	testAPI := mockPutObjectAPI(func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
		t.Helper()

		// Assertions to test the params given to the S3 client
		if *params.Key != testKey {
			t.Fatalf("expected supplied key to be %s but got %s", testKey, *params.Key)
		}
		if *params.Bucket != testBucket {
			t.Fatalf("expected supplied bucket to be %s but got %s", testBucket, *params.Bucket)
		}

		buf := new(strings.Builder)
		io.Copy(buf, params.Body)
		if buf.String() != testContents {
			t.Fatalf("unexpected file contents: '%s'", buf.String())
		}

		return &s3.PutObjectOutput{}, nil
	})

	err := saveToS3(ctx, testAPI, bytes.NewBuffer(contentsBuf), PutObjectParams{Bucket: testBucket, Key: testKey})
	if err != nil {
		t.Fatal(err)
	}
}
