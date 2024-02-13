package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
)

type GSPackage struct {
	Name     string
	Tag      string
	AssetUrl string
	Script   string
}

type GSConfig struct {
	GSPM     struct{}
	Packages []GSPackage
}

var args struct {
	Command string   `arg:"positional"`
	Values  []string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)

	var config GSConfig

	configFile, err1 := os.ReadFile("config.json")

	if err1 != nil {
		tui.ShowError("config.json not found")
	} else {
		err2 := json.Unmarshal(configFile, &config)

		if err2 != nil {
			tui.ShowError("Failed to parse config.json")
			os.Exit(1)
		}
	}

	tui.ShowInfo(fmt.Sprintf("Loaded %d packages from config.json", len(config.Packages)))

	if args.Command == "add" {
		tui.ShowInfo("Adding packages...")

		for _, value := range args.Values {
			tui.ShowLine()
			packageInfo := strings.Split(value, "@")

			var gsPackage GSPackage
			gsPackage.Name = packageInfo[0]

			if len(strings.Split(gsPackage.Name, "/")) != 2 {
				tui.ShowError("Invalid package name: " + gsPackage.Name)
				continue
			}

			if len(packageInfo) > 1 {
				gsPackage.Tag = packageInfo[1]
			}

			var needDownloadReleases bool = true

			for _, configPackage := range config.Packages {
				if configPackage.Name == gsPackage.Name {
					if gsPackage.Tag != "" && gsPackage.Tag != "latest" {
						if configPackage.Tag == gsPackage.Tag {
							needDownloadReleases = false
							gsPackage = configPackage
						}
					}
					break
				}
			}

			var res bool

			if needDownloadReleases {
				stopReleasesSpn := tui.ShowSpinner(fmt.Sprintf("Fetching %s releases...", gsPackage.Name))

				releases, err := gitservice.GetGitHubReleases(gsPackage.Name)

				if err != nil {
					stopReleasesSpn("fail")
					panic(err)
				}

				stopReleasesSpn("success")

				if len(releases) == 0 {
					tui.ShowInfo("No releases found for " + gsPackage.Name)
					continue
				}

				var tag string
				var tagOptions []string
				var assets = make(map[string][]gitservice.GSGitHubReleaseAsset)
				var assetOptions []string

				for _, release := range releases {
					if tag == "" {
						if gsPackage.Tag == "latest" {
							tag = release.TagName
						} else if gsPackage.Tag == release.TagName {
							tag = release.TagName
						}
						tagOptions = append(tagOptions, release.TagName)
					}
					assets[release.TagName] = release.Assets
				}

				if tag == "" {
					tag = tui.ShowOptions("Select a tag", tagOptions)
				}

				gsPackage.Tag = tag

				if len(assets[tag]) == 0 {
					tui.ShowInfo("No assets found for " + gsPackage.Name + "@" + tag)
					continue
				}

				for _, asset := range assets[tag] {
					assetOptions = append(assetOptions, asset.Name)
				}

				assetName := tui.ShowOptions("Select an asset", assetOptions)

				stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

				for _, asset := range assets[tag] {
					if asset.Name == assetName {
						gsPackage.AssetUrl = asset.BrowserDownloadUrl
						res = gitservice.GetGitHubReleaseAsset(asset.Name, asset.BrowserDownloadUrl)
						break
					}
				}

				if res {
					stopAssetSpn("success")

					runScript := tui.ShowConfirm("Do you want to run a script?")

					pkgPath := "PKG_PATH=~/Downloads/" + assetName

					if runScript {
						tui.ShowInfo(pkgPath)
						script := tui.ShowTextInput("Enter a script")
						gsPackage.Script = script
					}

					if len(gsPackage.Script) == 0 {
						tui.ShowInfo("Script not provided.")
						tui.ShowSuccess("Package located at ~/Downloads/" + assetName)
					} else {
						stopScriptSpn := tui.ShowSpinner("Running provided script...")

						tui.ShowBox(pkgPath + "\n" + gsPackage.Script)

						out, err := exec.Command("bash", "-c", pkgPath+"\n"+gsPackage.Script).Output()

						if err != nil {
							stopScriptSpn("fail")
							if string(out) != "" {
								tui.ShowBox(string(out))
							}
							tui.ShowError(err.Error())
						} else {
							stopScriptSpn("success")
							tui.ShowSuccess("Script executed successfully.")

							config.Packages = append(config.Packages, gsPackage)
						}
					}
				} else {
					stopAssetSpn("fail")
				}
			} else {
				tui.ShowInfo("Download Package from config")

				downloadUrlSplited := strings.Split(gsPackage.AssetUrl, "/")
				assetNameFromUrl := downloadUrlSplited[len(downloadUrlSplited)-1]

				stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetNameFromUrl))

				res = gitservice.GetGitHubReleaseAsset(assetNameFromUrl, gsPackage.AssetUrl)

				if res {
					stopAssetSpn("success")

					pkgPath := "PKG_PATH=~/Downloads/" + assetNameFromUrl

					if len(gsPackage.Script) == 0 {
						tui.ShowInfo("Script not provided.")
						tui.ShowSuccess("Package located at ~/Downloads/" + assetNameFromUrl)
					} else {
						stopScriptSpn := tui.ShowSpinner("Running provided script...")

						tui.ShowBox(pkgPath + "\n" + gsPackage.Script)

						out, err := exec.Command("bash", "-c", pkgPath+"\n"+gsPackage.Script).Output()

						if err != nil {
							stopScriptSpn("fail")
							if string(out) != "" {
								tui.ShowBox(string(out))
							}
							tui.ShowError(err.Error())
						} else {
							stopScriptSpn("success")
							tui.ShowSuccess("Script executed successfully.")
						}
					}
				} else {
					stopAssetSpn("fail")
				}
			}
		}

		// Save config as json file
		configBytes, _ := json.MarshalIndent(config, "", "    ")

		err := os.WriteFile("config.json", configBytes, 0644)

		if err != nil {
			panic(err)
		}
	} else {
		tui.ShowError("Unknown command: " + args.Command)
		os.Exit(1)
	}
}
