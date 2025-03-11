package gitservice

import (
	"fmt"

	"gitlab.com/gitlab-org/api/client-go"
)

func GitLabReleases(username string, repository string) ([]GSGitHubRelease, error) {
    client, err := gitlab.NewClient("")

	if err != nil {
		return nil, err
	}

	releases, _, err := client.Releases.ListReleases(fmt.Sprintf("%s/%s", username, repository), &gitlab.ListReleasesOptions{})

	var gsGitHubReleases []GSGitHubRelease

	for _, release := range releases {
		gsGitHubReleases = append(gsGitHubReleases, GSGitHubRelease{
			//Url:             *release.URL,
			//AssetsUrl:       *release.AssetsURL,
			//UploadUrl:       *release.UploadURL,
			//HtmlUrl:         *release.HTMLURL,
			//Id:              *release.ID,
			TagName:         release.TagName,
			//TargetCommitish: *release.TargetCommitish,
			//Name:            *release.Name,
			//Draft:           *release.Draft,
			//Prerelease:      *release.Prerelease,
			//CreatedAt:       *release.CreatedAt,
			//PublishedAt:     *release.PublishedAt,
		})
	}

	return gsGitHubReleases, nil
}

func GitLabReleaseAssets(username string, repository string, tagName string) ([]GSGitHubReleaseAsset, error) {
    client, err := gitlab.NewClient("")

	if err != nil {
		return nil, err
	}

	releases, _, err := client.Releases.ListReleases(fmt.Sprintf("%s/%s", username, repository), &gitlab.ListReleasesOptions{})

	var gsGitHubReleaseAssets []GSGitHubReleaseAsset

	for _, release := range releases {
        if release.TagName == tagName {
            for _, link := range release.Assets.Links {
			    gsGitHubReleaseAssets = append(gsGitHubReleaseAssets, GSGitHubReleaseAsset{
			        Url:                link.URL,
			        //Id:                 *asset.ID,
			        //NodeId:             *asset.NodeID,
			        Name:               link.Name,
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
			break
		}
	}

	return gsGitHubReleaseAssets, nil
}

func GitLabReleaseAssetDownload(url string, name string) (bool, error) {
	res, err := GetGitHubReleaseAsset(name, url)
	return res, err
}
