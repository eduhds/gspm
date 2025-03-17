package gitservice

import (
	"fmt"
)

func BitbucketReleases(username string, repository string) ([]GSRelease, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/downloads", username, repository)

	releases, err := GSGetReleases(url, "")

	if err != nil {
		return nil, err
	}

	var gsReleases []GSRelease

	releasesMap := releases.(map[string]any)
	values := releasesMap["values"].([]any)

	for _, release := range values {
		releaseMap := release.(map[string]any)
		gsReleases = append(gsReleases, GSRelease{
			TagName: releaseMap["name"].(string),
		})
	}

	return gsReleases, nil
}

func BitbucketReleaseAssets(username string, repository string, tagName string) ([]GSReleaseAsset, error) {
	var gsReleaseAssets []GSReleaseAsset

	return gsReleaseAssets, nil
}

func BitbucketReleaseAssetDownload(url string, name string) (bool, error) {
	return res, err
}

