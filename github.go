package main

import (
	"context"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

func githubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *inputRepoTokenPtr},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}
