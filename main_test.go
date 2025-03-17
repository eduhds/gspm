package main

import (
	"runtime"
	"strings"
	"testing"

	"github.com/eduhds/gspm/internal/gitservice"
)

const (
	username   string = "eduhds"
	repository string = "gspm"
)

/*
 * GITHUB TESTS
 */

 func TestGitHubReleases(t *testing.T) {
	releases, err := gitservice.GitHubReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	if len(releases) == 0 {
		t.Fatal("No releases found")
	}

	if releases[0].TagName != "v"+version {
		t.Fatal("No release found for version " + version)
	}
}

func TestGitHubReleaseAssets(t *testing.T) {
	releases, err := gitservice.GitHubReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	assets, err := gitservice.GitHubReleaseAssets(username, repository, releases[len(releases)-1].Id)

	if err != nil {
		t.Fatal(err)
	}

	if len(assets) == 0 {
		t.Fatal("No assets found")
	}
}

func TestGitHubDownloadAsset(t *testing.T) {
	releases, err := gitservice.GitHubReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	assets, err := gitservice.GitHubReleaseAssets(username, repository, releases[0].Id)

	if err != nil {
		t.Fatal(err)
	}

	var (
		assetName string
		assetUrl  string
	)

	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), runtime.GOOS) && strings.Contains(strings.ToLower(asset.Name), version) {
			assetName = asset.Name
			assetUrl = asset.Url
			break
		}
	}

	downlod, err := gitservice.GitHubReleaseAssetDownload(assetUrl, assetName)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}

/*
 * GITLAB TESTS
 */

func TestGitLabReleases(t *testing.T) {
	releases, err := gitservice.GitLabReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	if len(releases) == 0 {
		t.Fatal("No releases found")
	}
}

func TestGitLabReleaseAssets(t *testing.T) {
	releases, err := gitservice.GitLabReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	assets, err := gitservice.GitLabReleaseAssets(username, repository, releases[0].TagName)

	if err != nil {
		t.Fatal(err)
	}

	if len(assets) == 0 {
		t.Fatal("No assets found")
	}
}

func TestGitLabDownloadAsset(t *testing.T) {
	releases, err := gitservice.GitLabReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	assets, err := gitservice.GitLabReleaseAssets(username, repository, releases[0].TagName)

	if err != nil {
		t.Fatal(err)
	}

	var (
		assetName string
		assetUrl  string
	)

	for _, asset := range assets {
		if strings.Contains(strings.ToLower(asset.Name), runtime.GOOS) && strings.Contains(strings.ToLower(asset.Name), version) {
			assetName = asset.Name
			assetUrl = asset.Url
			break
		}
	}

	downlod, err := gitservice.GitLabReleaseAssetDownload(assetUrl, assetName)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}

/*
 * BITBUCKET TESTS
 */

func TestBitbucketReleases(t *testing.T) {
	releases, err := gitservice.BitbucketReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	if len(releases) == 0 {
		t.Fatal("No releases found")
	}
}

