package main

import (
	"fmt"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
)

type GSPackage struct {
	Name     string
	Tag      string
	AssetUrl string
}

var args struct {
	Command string `arg:"positional"`
	Value   string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)

	if args.Command == "add" {
		packageInfo := strings.Split(args.Value, "@")

		var gsPackage GSPackage
		gsPackage.Name = packageInfo[0]

		if len(strings.Split(gsPackage.Name, "/")) != 2 {
			panic("Invalid package name")
		}

		if len(packageInfo) > 1 {
			gsPackage.Tag = packageInfo[1]
		}

		stopReleasesSpn := tui.ShowSpinner(fmt.Sprintf("Fetching %s releases...", gsPackage.Name))

		releases, err := gitservice.GetGitHubReleases(gsPackage.Name)

		if err != nil {
			stopReleasesSpn("fail")
			panic(err)
		}

		stopReleasesSpn("success")

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
			tag = tui.ShowOptions("Select a tag", tagOptions)
		}

		gsPackage.Tag = tag

		for _, asset := range assets[tag] {
			assetOptions = append(assetOptions, asset.Name)
		}

		assetName := tui.ShowOptions("Select an asset", assetOptions)

		stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

		var res bool

		for _, asset := range assets[tag] {
			if asset.Name == assetName {
				res = gitservice.GetGitHubReleaseAsset(asset)
				break
			}
		}

		if res {
			stopAssetSpn("success")
		} else {
			stopAssetSpn("fail")
		}
	} else {
		panic("Unknown command: " + args.Command)
	}
}
