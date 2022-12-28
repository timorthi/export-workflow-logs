package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestGetWorkflowRunLogsURLForRunID(t *testing.T) {
	t.Setenv(envVarRepoOwner, "someowner")
	t.Setenv(envVarRepoFullName, "someowner/test-repo")
	testURL := "https://example.com/some/test/path/to/logs.zip"

	mockedHTTPClient := mock.NewMockedHTTPClient(
		mock.WithRequestMatchHandler(
			mock.GetReposActionsRunsLogsByOwnerByRepoByRunId,
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Location", testURL)
				w.WriteHeader(http.StatusFound)
			}),
		),
	)
	testClient := github.NewClient(mockedHTTPClient)

	url, err := getWorkflowRunLogsURLForRunID(context.Background(), testClient, 123)

	if err != nil {
		t.Fatal(err)
	}

	if url.String() != testURL {
		t.Fatalf("Expected url to equal '%s' but got: '%s'", testURL, url.String())
	}
}
