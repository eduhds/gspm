package gitservice

import (
	"context"

	"github.com/google/go-github/v69/github"
)

var GHToken = ""

func GitHubReleases(username string, repository string) ([]GSRelease, error) {
	client := github.NewClient(nil)

	releases, _, err := client.Repositories.ListReleases(context.Background(), username, repository, nil)

	var gsReleases []GSRelease

	for _, release := range releases {
		gsReleases = append(gsReleases, GSRelease{
			Id:      *release.ID,
			TagName: *release.TagName,
		})
	}

	return gsReleases, err
}

func GitHubReleaseAssets(username string, repository string, id int64) ([]GSReleaseAsset, error) {
	client := github.NewClient(nil)

	assets, _, err := client.Repositories.ListReleaseAssets(context.Background(), username, repository, id, nil)

	var gsReleaseAssets []GSReleaseAsset

	for _, asset := range assets {
		gsReleaseAssets = append(gsReleaseAssets, GSReleaseAsset{
			Url:  *asset.URL,
			Id:   *asset.ID,
			Name: *asset.Name,
		})
	}

	return gsReleaseAssets, err
}

func GitHubReleaseAssetDownload( /* username string, repository string, id int64, */ url string, name string) (bool, error) {
	/* client := github.NewClient(nil)

	_, url, err := client.Repositories.DownloadReleaseAsset(context.Background(), username, repository, id, nil)

	if err != nil {
		return false, err
	} */

	res, err := GSGetReleaseAsset(name, url, GHToken)

	return res, err
}
