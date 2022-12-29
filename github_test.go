package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-github/v48/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	assert.Equal(t, url.String(), testURL)
}
