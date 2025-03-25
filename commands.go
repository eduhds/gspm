package main

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
	"github.com/eduhds/gspm/internal/util"
)

var serviceSymbol = map[string]string{
	"gitlab": "🦊",
	"github": "🐙",
	"bitbucket": "🪣 ",
}

func CommandAdd(cfg GSConfig, gsp GSPackage) GSConfig {
	stopReleasesSpn := tui.ShowSpinner(fmt.Sprintf("Fetching %s releases...", gsp.Name))

	releases, err := ResolveReleases(gsp)

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

	var tag string
	var tagOptions []string
	var assets = make(map[string][]gitservice.GSReleaseAsset)
	var assetOptions []string

	for _, release := range releases {
		//if tag == "" {
		//	if gsp.Tag == release.TagName {
		//		tag = release.TagName
		//	}
		tagOptions = append(tagOptions, release.TagName)
		//}
		assets[release.TagName] = release.Assets
	}

	//if tag == "" {
	tag = tui.ShowOptions("Select a tag", tagOptions)
	gsp.Tag = tag
	//}

	if len(assets[tag]) == 0 {
		tui.ShowInfo("No assets found for " + gsp.Name + "@" + tag)
		return cfg
	}

	for _, asset := range assets[tag] {
		assetOptions = append(assetOptions, asset.Name)
	}

	assetName := tui.ShowOptions("Select an asset", assetOptions)

	for _, asset := range assets[tag] {
		if asset.Name == assetName {
			gsp.AssetUrl = asset.Url
			break
		}
	}

	stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

	success, err := ResolveAsset(gsp.AssetUrl, assetName)
	if err != nil {
		tui.ShowError(err.Error())
	}

	if success {
		stopAssetSpn("success")

		if gsp.Script == "" {
			wantRunScript := tui.ShowConfirm("Do you want to run a script?")

			if wantRunScript {
				tui.ShowInfo("Use {{ASSET}} to reference the asset path")
				gsp.Script = tui.ShowTextInput("Enter a script", true, "")
			} else {
				tui.ShowInfo("Script not provided.")
				tui.ShowSuccess("Package located at " + filepath.Join(downloadPrefix, assetName))
			}
		}

		if gsp.Script != "" {
			if !RunScript(assetName, gsp.Script) {
				return cfg
			}
		}

		gsp.LastModfied = util.DateLocalNow()
		gsp.Service = service

		var found = false

		for index, configPackage := range cfg.Packages {
			if configPackage.Platform == runtime.GOOS && configPackage.Name == gsp.Name && ((configPackage.Service != "" && configPackage.Service == service) || (service == gitservice.GITHUB)) {
				cfg.Packages[index] = gsp
				found = true
				break
			}
		}

		if !found {
			cfg.Packages = append(cfg.Packages, gsp)
		}
	} else {
		stopAssetSpn("fail")
	}

	return cfg
}

func CommandEdit(cfg GSConfig, gsp GSPackage, withScript bool) GSConfig {
	tui.ShowWarning("Editing script for package " + gsp.Name)

	if withScript {
		tui.ShowMessage(gsp.Script)
	} else {
		tui.ShowInfo("Use {{ASSET}} to reference the asset path")
		script := tui.ShowTextInput("Enter a script", true, gsp.Script)
		gsp.Script = script
	}

	for index, configPackage := range cfg.Packages {
		if configPackage.Platform == runtime.GOOS && configPackage.Name == gsp.Name {
			gsp.LastModfied = util.DateLocalNow()
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
		fmt.Println("Service: " + gsp.Service)
		fmt.Println("LastModfied: " + gsp.LastModfied)
	}
}

func CommandInstall(config GSConfig) {
	platformPackages := PlatformPackages(config)

	if len(platformPackages) == 0 {
		tui.ShowInfo("No packages found")
		return
	}

	tui.ShowInfo(fmt.Sprintf("Loaded %d packages", len(platformPackages)))

	for _, item := range platformPackages {
		tui.ShowInfo("Installing package: " + item.Name)
		assetName := AssetNameFromUrl(item.AssetUrl)

		stopInstallSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

		success, err := ResolveAsset(item.AssetUrl, assetName)
		if err != nil {
			tui.ShowError(err.Error())
		}

		if success {
			stopInstallSpn("success")
			RunScript(assetName, item.Script)
		} else {
			stopInstallSpn("fail")
		}
	}
}

func CommandList(config GSConfig) {
	platformPackages := PlatformPackages(config)

	if len(platformPackages) == 0 {
		tui.ShowInfo("No packages found")
		return
	}

	tui.ShowInfo(fmt.Sprintf("%d packages for %s", len(platformPackages), runtime.GOOS))
	tui.ShowLine()

	for _, item := range platformPackages {
		tui.TextInfo("📦 " + item.Name + "@" + item.Tag)
		tui.ShowMessage("🔗 " + item.AssetUrl)
		tui.ShowMessage("🛠️  " + item.Script)
		if item.Service != "" {
			tui.ShowMessage(serviceSymbol[item.Service] + " " + item.Service)
		}
		tui.ShowLine()
	}
}

func CommandRemove(cfg GSConfig, gsp GSPackage, withScript bool) GSConfig {
	var keepedPackages []GSPackage

	for _, item := range cfg.Packages {
		if item.Platform == runtime.GOOS && item.Name == gsp.Name {
			if withScript {
				if !RunScript("", gsp.Script) {
					return cfg
				}
				tui.ShowInfo("Package " + item.Name + " removed with script")
			} else {
				tui.ShowInfo("Package " + item.Name + " removed from config")
			}
		} else {
			keepedPackages = append(keepedPackages, item)
		}
	}

	cfg.Packages = keepedPackages

	return cfg
}
