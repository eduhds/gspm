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

func IsBestPackageToDownload(assetName string) bool {
	assetName = strings.ToLower(assetName)
	isSameOS := strings.Contains(assetName, runtime.GOOS)
	isSameArch := strings.Contains(assetName, "x86_64")
	isSameVersion := strings.Contains(assetName, version)
	isTarGz := strings.HasSuffix(assetName, ".tar.gz")

	return isSameOS && isSameArch && isSameVersion && isTarGz
}

//
// GITHUB TESTS
//

func TestGitHubDownloadAsset(t *testing.T) {
	releases, err := gitservice.GitHubReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	var (
		assetName string
		assetUrl  string
	)

	for _, asset := range releases[0].Assets {
		if IsBestPackageToDownload(asset.Name) {
			assetName = asset.Name
			assetUrl = asset.Url
			break
		}
	}

	downlod, err := gitservice.GitHubReleaseAssetDownload(assetUrl, "TEST_GITHUB_"+assetName)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}

//
// GITLAB TESTS
//

func TestGitLabDownloadAsset(t *testing.T) {
	releases, err := gitservice.GitLabReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	var (
		assetName string
		assetUrl  string
	)

	for _, asset := range releases[0].Assets {
		if IsBestPackageToDownload(asset.Name) {
			assetName = asset.Name
			assetUrl = asset.Url
			break
		}
	}

	downlod, err := gitservice.GitLabReleaseAssetDownload(assetUrl, "TEST_GITLAB_"+assetName)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}

//
// BITBUCKET TESTS
//

func TestBitbucketDownloadAsset(t *testing.T) {
	releases, err := gitservice.BitbucketReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

	var (
		assetName string
		assetUrl  string
	)

	for _, asset := range releases[0].Assets {
		if IsBestPackageToDownload(asset.Name) {
			assetName = asset.Name
			assetUrl = asset.Url
			break
		}
	}

	downlod, err := gitservice.BitbucketReleaseAssetDownload(assetUrl, "TEST_BITBUCKET_"+assetName)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}
