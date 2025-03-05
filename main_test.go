package main

import (
	"fmt"
	"testing"

	"github.com/eduhds/gspm/internal/gitservice"
)

func TestGitHubReleases(t *testing.T) {
	username := "eduhds"
	repository := "gspm"
	releases, err := gitservice.GitHubReleases(username, repository)
	if err != nil {
		t.Fatal(err)
	}
fmt.Println(releases)
}

func TestGitHubReleaseAssets(t *testing.T) {
	username := "eduhds"
	repository := "gspm"
	releases, err := gitservice.GitHubReleases(username, repository)

	assets, err := gitservice.GitHubReleaseAssets(username, repository, releases[len(releases)-1].Id)

	if err != nil {
		t.Fatal(err)
	}

fmt.Println(assets)
}

