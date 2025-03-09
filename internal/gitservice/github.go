package gitservice

import (
	"context"

	"github.com/google/go-github/v69/github"
)

func GitHubReleases(username string, repository string) ([]GSGitHubRelease, error) {
	client := github.NewClient(nil)

	releases, _, err := client.Repositories.ListReleases(context.Background(), username, repository, nil)

	var gsGitHubReleases []GSGitHubRelease

	for _, release := range releases {
		gsGitHubReleases = append(gsGitHubReleases, GSGitHubRelease{
			Url:             *release.URL,
			AssetsUrl:       *release.AssetsURL,
			UploadUrl:       *release.UploadURL,
			HtmlUrl:         *release.HTMLURL,
			Id:              *release.ID,
			TagName:         *release.TagName,
			TargetCommitish: *release.TargetCommitish,
			Name:            *release.Name,
			Draft:           *release.Draft,
			Prerelease:      *release.Prerelease,
			//CreatedAt:       *release.CreatedAt,
			//PublishedAt:     *release.PublishedAt,
		})
	}

	return gsGitHubReleases, err
}

func GitHubReleaseAssets(username string, repository string, id int64) ([]GSGitHubReleaseAsset, error) {
	client := github.NewClient(nil)

	assets, _, err := client.Repositories.ListReleaseAssets(context.Background(), username, repository, id, nil)

	var gsGitHubReleaseAssets []GSGitHubReleaseAsset

	for _, asset := range assets {
		gsGitHubReleaseAssets = append(gsGitHubReleaseAssets, GSGitHubReleaseAsset{
			Url:                *asset.URL,
			Id:                 *asset.ID,
			//NodeId:             *asset.NodeID,
			Name:               *asset.Name,
			//Label:              *asset.Label,
			//ContentType:        *asset.ContentType,
			//State:              *asset.State,
			//Size:               *asset.Size,
			//DownloadCount:      *asset.DownloadCount,
			//CreatedAt:          *asset.CreatedAt,
			//UpdatedAt:          *asset.UpdatedAt,
			//BrowserDownloadUrl: *asset.BrowserDownloadURL,
		})
	}

	return gsGitHubReleaseAssets, err
}

func GitHubReleaseAssetDownload(username string, repository string, id int64, name string) (bool, error) {
	client := github.NewClient(nil)

	_, url, err := client.Repositories.DownloadReleaseAsset(context.Background(), username, repository, id, nil)

	if err != nil {
		return false, err
	}

    res, err := GetGitHubReleaseAsset(name, url)

	return res, err
}

