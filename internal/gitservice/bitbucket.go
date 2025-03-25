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
		var gsReleaseAssets []GSReleaseAsset

		gsReleaseAssets = append(gsReleaseAssets, GSReleaseAsset{
			Url:  release.Links.Self.Href,
			Name: release.Name,
		})

		// Obs.: Source tarballs must be fetched separately

		gsReleases = append(gsReleases, GSRelease{
			TagName: release.Name,
			Assets:  gsReleaseAssets,
		})
	}

	return gsReleases, nil
}

func BitbucketReleaseAssetDownload(url string, name string) (bool, error) {
	res, err := GSGetReleaseAsset(name, url, BBToken)
	return res, err
}
