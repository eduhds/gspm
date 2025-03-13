package gitservice

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GitLabReleases(username string, repository string) ([]GSRelease, error) {
	client, err := gitlab.NewClient("")

	if err != nil {
		return nil, err
	}

	releases, _, err := client.Releases.ListReleases(fmt.Sprintf("%s/%s", username, repository), &gitlab.ListReleasesOptions{})

	var gsReleases []GSRelease

	for _, release := range releases {
		gsReleases = append(gsReleases, GSRelease{
			TagName: release.TagName,
		})
	}

	return gsReleases, nil
}

func GitLabReleaseAssets(username string, repository string, tagName string) ([]GSReleaseAsset, error) {
	client, err := gitlab.NewClient("")

	if err != nil {
		return nil, err
	}

	releases, _, err := client.Releases.ListReleases(fmt.Sprintf("%s/%s", username, repository), &gitlab.ListReleasesOptions{})

	var gsReleaseAssets []GSReleaseAsset

	for _, release := range releases {
		if release.TagName == tagName {
			for _, link := range release.Assets.Links {
				gsReleaseAssets = append(gsReleaseAssets, GSReleaseAsset{
					Url:  link.URL,
					Name: link.Name,
				})
			}
			break
		}
	}

	return gsReleaseAssets, nil
}

func GitLabReleaseAssetDownload(url string, name string) (bool, error) {
	res, err := GetGitHubReleaseAsset(name, url)
	return res, err
}
