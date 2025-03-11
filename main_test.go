package main

import (
	"testing"

	"github.com/eduhds/gspm/internal/gitservice"
)

const (
	username string = "eduhds"
	repository string = "gspm"
)

func TestGitHubReleases(t *testing.T) {
	releases, err := gitservice.GitHubReleases(username, repository)

	if err != nil {
		t.Fatal(err)
	}

    if len(releases) == 0 {
		t.Fatal("No releases found")
	}

	if releases[0].TagName != "v" + version {
		t.Fatal("No release found for version " + version)
	}
}

func TestGitHubReleaseAssets(t *testing.T) {
	releases, err := gitservice.GitHubReleases(username, repository)

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

	assets, err := gitservice.GitHubReleaseAssets(username, repository, releases[0].Id)

	downlod, err := gitservice.GitHubReleaseAssetDownload(username, repository, assets[0].Id, assets[0].Name)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}

func TestGitLabReleases(t *testing.T) {
	releases, err := gitservice.GitLabReleases("Inkscape", "inkscape")

	if err != nil {
		t.Fatal(err)
	}

	if len(releases) == 0 {
		t.Fatal("No releases found")
	}
}

func TestGitLabReleaseAssets(t *testing.T) {
	releases, err := gitservice.GitLabReleases("Inkscape", "inkscape")

	assets, err := gitservice.GitLabReleaseAssets("Inkscape", "inkscape", releases[0].TagName)

	if err != nil {
		t.Fatal(err)
	}

	if len(assets) == 0 {
		t.Fatal("No assets found")
	}
}

func TestGitLabDownloadAsset(t *testing.T) {
	releases, err := gitservice.GitLabReleases("Inkscape", "inkscape")

	assets, err := gitservice.GitLabReleaseAssets("Inkscape", "inkscape", releases[0].TagName)

	downlod, err := gitservice.GitLabReleaseAssetDownload(assets[0].Url, assets[0].Name)

	if err != nil {
		t.Fatal(err)
	}

	if !downlod {
		t.Fatal("Asset not downloaded")
	}
}

