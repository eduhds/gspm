package gitservice

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

var GLToken = ""

func GitLabReleases(username string, repository string) ([]GSRelease, error) {
	client, err := gitlab.NewClient(GLToken)

	if err != nil {
		return nil, err
	}

	releases, _, err := client.Releases.ListReleases(fmt.Sprintf("%s/%s", username, repository), &gitlab.ListReleasesOptions{})

	var gsReleases []GSRelease

	for _, release := range releases {
		var gsReleaseAssets []GSReleaseAsset

		for _, link := range release.Assets.Links {
			gsReleaseAssets = append(gsReleaseAssets, GSReleaseAsset{
				Url:  link.URL,
				Name: link.Name,
			})
		}

		for _, source := range release.Assets.Sources {
			gsReleaseAssets = append(gsReleaseAssets, GSReleaseAsset{
				Url:  source.URL,
				Name: fmt.Sprintf("%s.%s", repository, source.Format),
			})
		}

		gsReleases = append(gsReleases, GSRelease{
			TagName: release.TagName,
			Assets:  gsReleaseAssets,
		})
	}

	return gsReleases, nil
}

func GitLabReleaseAssetDownload(url string, name string) (bool, error) {
	res, err := GSGetReleaseAsset(name, url, GLToken)
	return res, err
}
