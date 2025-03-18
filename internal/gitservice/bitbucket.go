package gitservice

import (
	"fmt"
)

var BBToken = ""

func BitbucketReleases(username string, repository string) ([]GSRelease, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/downloads", username, repository)

	var releases BitbucketResponse
	_, err := GSGetReleases(url, BBToken, &releases)

	if err != nil {
		return nil, err
	}

	var gsReleases []GSRelease

	for _, release := range releases.Values {
		gsReleases = append(gsReleases, GSRelease{
			TagName: release.Name,
		})
	}

	return gsReleases, nil
}

func BitbucketReleaseAssets(username string, repository string, tagName string) ([]GSReleaseAsset, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/downloads", username, repository)

	var releases BitbucketResponse
	_, err := GSGetReleases(url, BBToken, &releases)

	if err != nil {
		return nil, err
	}

	var gsReleaseAssets []GSReleaseAsset

	for _, release := range releases.Values {
		if release.Name == tagName {
			gsReleaseAssets = append(gsReleaseAssets, GSReleaseAsset{
				Url:  release.Links.Self.Href,
				Name: release.Name,
			})
		}
	}

	return gsReleaseAssets, nil
}

func BitbucketReleaseAssetDownload(url string, name string) (bool, error) {
	res, err := GSGetReleaseAsset(name, url, BBToken)
	return res, err
}
