package main

import (
	"context"
	"io"
	"os"
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
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "tmpFile")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()
	testContents := "hello world!"
	tmpFile.WriteString(testContents)

	testBucket := "testBucketName"
	testKey := "some/key/to/logs.zip"

	testAPI := mockPutObjectAPI(func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
		t.Helper()
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

	err = saveToS3(ctx, testAPI, testBucket, testKey, tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Fatalf("expected temp dir to have been cleaned up, but got err: %v", err)
	}
}
