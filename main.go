package main

import (
	"fmt"

	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
)

type GSPackage struct {
	Name     string
	Tag      string
	AssetUrl string
}

func main() {
	tui.ShowMessage("Welcome to GSPM!")

	var gsPackage GSPackage
	gsPackage.Name = "eduhds/logduto"
	gsPackage.Tag = "latest"

	stopSpinner := tui.ShowSpinner("Fetching releases...")

	releases, err := gitservice.GetGitHubReleases(gsPackage.Name)

	if err != nil {
		stopSpinner("fail")
		panic(err)
	}

	stopSpinner("success")

	var tag string
	var tagOptions []string
	var assets = make(map[string][]gitservice.GSGitHubReleaseAsset)
	var assetOptions []string

	for _, release := range releases {
		if gsPackage.Tag != "" && release.TagName == gsPackage.Tag {
			tag = release.TagName
			break
		}
		tagOptions = append(tagOptions, release.TagName)
		assets[release.TagName] = release.Assets
	}

	if tag == "" {
		tag = tui.ShowOptions("\nSelect a tag", tagOptions)
	}

	gsPackage.Tag = tag

	for _, asset := range assets[tag] {
		assetOptions = append(assetOptions, asset.Name)
	}

	assetName := tui.ShowOptions("\nSelect an asset", assetOptions)

	var res bool

	for _, asset := range assets[tag] {
		if asset.Name == assetName {
			res = gitservice.GetGitHubReleaseAsset(asset)
			break
		}
	}

	tui.ShowInfo(fmt.Sprintf("Result %v", res))
}
