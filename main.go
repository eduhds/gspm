package main

import (
	"fmt"

	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
)

type Package struct {
	Name     string
	Tag      string
	AssetUrl string
}

func main() {
	fmt.Println("hello, world")

	releases, err := gitservice.GetGitHubReleases("eduhds", "logduto")

	if err != nil {
		panic(err)
	}

	var options []string
	var assets = make(map[string][]gitservice.GSGitHubReleaseAsset)

	for _, release := range releases {
		options = append(options, release.TagName)
		assets[release.TagName] = release.Assets
	}

	opt := tui.ShowOptions(options)
	tui.ShowMessage(opt)

	options = []string{}

	for _, asset := range assets[opt] {
		options = append(options, asset.Name)
	}

	file := tui.ShowOptions(options)
	tui.ShowMessage(file)

	var res bool

	for _, asset := range assets[opt] {
		if asset.Name == file {
			fmt.Println(asset.BrowserDownloadUrl)
			res = gitservice.GetGitHubReleaseAsset(asset)
			break
		}
	}

	tui.ShowMessage(fmt.Sprintf("%v", res))
}
