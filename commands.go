package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
)

func CommandEdit(cfg GSConfig, gsp GSPackage) GSConfig {
	tui.ShowWarning("Editing script for package " + gsp.Name)

	tui.ShowInfo("Use {{ASSET}} to reference the asset path")
	script := tui.ShowTextInput("Enter a script", true, gsp.Script)
	gsp.Script = script

	for index, configPackage := range cfg.Packages {
		if configPackage.Platform == runtime.GOOS && configPackage.Name == gsp.Name {
			cfg.Packages[index] = gsp
			break
		}
	}

	return cfg
}

func CommandInfo(gsp GSPackage) {
	if gsp.Tag != "" {
		fmt.Println("Name: " + gsp.Name)
		fmt.Println("Tag: " + gsp.Tag)
		fmt.Println("AssetUrl: " + gsp.AssetUrl)
		fmt.Println("Script: " + gsp.Script)
		fmt.Println("Platform: " + gsp.Platform)
		fmt.Println("LastModfied: " + gsp.LastModfied)
	}
}

func CommandInstall(platformPackages []GSPackage) {
	tui.ShowInfo(fmt.Sprintf("Loaded %d packages", len(platformPackages)))

	for _, item := range platformPackages {
		tui.ShowInfo("Installing package: " + item.Name)
		assetName := AssetNameFromUrl(item.AssetUrl)

		stopInstallSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

		success := gitservice.GetGitHubReleaseAsset(assetName, item.AssetUrl)

		if success {
			stopInstallSpn("success")
			RunScript(assetName, item.Script)
		} else {
			stopInstallSpn("fail")
		}
	}
}

func CommandList(platformPackages []GSPackage) {
	tui.ShowInfo(fmt.Sprintf("%d packages for %s", len(platformPackages), runtime.GOOS))
	tui.ShowLine()

	for _, item := range platformPackages {
		tui.ShowMessage("📦 " + item.Name + "@" + item.Tag)
		tui.ShowMessage("🔗 " + item.AssetUrl)
		tui.ShowMessage("🛠️  " + item.Script)
		tui.ShowLine()
	}
}

func CommandRemove(cfg GSConfig, gsp GSPackage) GSConfig {
	var keepedPackages []GSPackage

	for _, item := range cfg.Packages {
		if item.Platform == runtime.GOOS && item.Name == gsp.Name {
			tui.ShowInfo("Package " + item.Name + " removed")
		} else {
			keepedPackages = append(keepedPackages, item)
		}
	}

	cfg.Packages = keepedPackages

	return cfg
}

func CommandUpdate(cfg GSConfig, gsp GSPackage) GSConfig {

	var assetName = ""

	assetName = AssetNameFromUrl(gsp.AssetUrl)

	stopReleasesSpn := tui.ShowSpinner(fmt.Sprintf("Fetching %s releases...", gsp.Name))

	releases, err := gitservice.GetGitHubReleases(gsp.Name)

	if err != nil {
		stopReleasesSpn("fail")
		tui.ShowError(err.Error())
		return cfg
	}

	stopReleasesSpn("success")

	if len(releases) == 0 {
		tui.ShowInfo("No releases found for " + gsp.Name)
		return cfg
	}

	if releases[0].TagName == gsp.Tag {
		tui.ShowInfo("No updates found for " + gsp.Name + ", already at latest tag " + gsp.Tag)
		return cfg
	}

	assets := releases[0].Assets
	gsp.Tag = releases[0].TagName

	if len(assets) == 0 {
		tui.ShowInfo("No assets found for " + gsp.Name + "@" + gsp.Tag)
		return cfg
	}

	var assetOptions []string

	for _, asset := range assets {
		fmt.Println(asset.Name)
		if assetName == asset.Name {
			gsp.AssetUrl = asset.BrowserDownloadUrl
			break
		}
		assetOptions = append(assetOptions, asset.Name)
	}

	if len(assetOptions) == len(assets) {
		tui.ShowError("Asset not found: " + assetName + ". Need select another.")

		assetName = tui.ShowOptions("Select an asset", assetOptions)

		for _, asset := range assets {
			if asset.Name == assetName {
				gsp.AssetUrl = asset.BrowserDownloadUrl
				break
			}
		}
	}

	stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

	success := gitservice.GetGitHubReleaseAsset(assetName, gsp.AssetUrl)

	if success {
		stopAssetSpn("success")

		if RunScript(assetName, gsp.Script) {
			gsp.LastModfied = time.Now().Local().Format(time.RFC3339)

			for index, configPackage := range cfg.Packages {
				if configPackage.Platform == runtime.GOOS && configPackage.Name == gsp.Name {
					cfg.Packages[index] = gsp
					break
				}
			}
		}
	} else {
		stopAssetSpn("fail")
	}

	return cfg
}
