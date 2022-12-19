package main

import (
	"context"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func GithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GetRequiredEnv("INPUT_REPO-TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}
