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

func RunScript(assetPath string, providedScript string) bool {
	script := "ASSET_PATH=" + assetPath + "\n" + providedScript

	stopScriptSpn := tui.ShowSpinner("Running provided script...")

	tui.ShowBox(script)

	out, err := exec.Command("bash", "-c", script).Output()

	if err != nil {
		stopScriptSpn("fail")
		if string(out) != "" {
			tui.ShowBox(string(out))
		}
		tui.ShowError(err.Error())
		return false
	} else {
		stopScriptSpn("success")
		tui.ShowSuccess("Script executed successfully.")
		return true
	}
}

func ReadConfig() GSConfig {
	var config GSConfig

	configFile, err := os.ReadFile("config.json")

	if err != nil {
		tui.ShowError("config.json not found")
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			tui.ShowError("Failed to parse config.json")
		}
	}

	return config
}

func main() {
	arg.MustParse(&args)

	config := ReadConfig()
	countPackages := len(config.Packages)

	if args.Command == "install" {
		if countPackages == 0 {
			tui.ShowInfo("No packages found in config.json")
			return
		}
		tui.ShowInfo(fmt.Sprintf("Loaded %d packages from config.json", countPackages))
		// TODO: Loop through packages and install them
	} else if args.Command == "add" {
		if len(args.Values) == 0 {
			tui.ShowError("No packages provided")
			return
		}

		tui.ShowInfo("Adding packages...")

		for _, value := range args.Values {
			tui.ShowLine()
			packageInfo := strings.Split(value, "@")
			packageName := packageInfo[0]
			packageTag := ""

			if len(strings.Split(packageName, "/")) != 2 {
				tui.ShowError("Invalid package name: " + packageName)
				continue
			}

			if len(packageInfo) > 1 {
				packageTag = packageInfo[1]
			}

			var gsPackage GSPackage
			gsPackage.Name = packageName
			gsPackage.Tag = packageTag

			var needDownloadReleases bool = true

			if gsPackage.Tag != "latest" {
				for _, configPackage := range config.Packages {
					if configPackage.Name == gsPackage.Name {
						if gsPackage.Tag == "" || configPackage.Tag == gsPackage.Tag {
							needDownloadReleases = false
							gsPackage = configPackage
						}
						break
					}
				}
			}

			var res bool

			if needDownloadReleases {
				stopReleasesSpn := tui.ShowSpinner(fmt.Sprintf("Fetching %s releases...", gsPackage.Name))

				releases, err := gitservice.GetGitHubReleases(gsPackage.Name)

				if err != nil {
					stopReleasesSpn("fail")
					tui.ShowError(err.Error())
					return
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

					if runScript {
						tui.ShowInfo("Tip: Use $ASSET_PATH to reference the asset path")
						script := tui.ShowTextInput("Enter a script")
						gsPackage.Script = script
					}

					if len(gsPackage.Script) == 0 {
						tui.ShowInfo("Script not provided.")
						tui.ShowSuccess("Package located at ~/Downloads/" + assetName)
					} else {
						if RunScript("~/Downloads/"+assetName, gsPackage.Script) {
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

					if len(gsPackage.Script) == 0 {
						tui.ShowInfo("Script not provided.")
						tui.ShowSuccess("Package located at ~/Downloads/" + assetNameFromUrl)
					} else {
						RunScript("~/Downloads/"+assetNameFromUrl, gsPackage.Script)
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
