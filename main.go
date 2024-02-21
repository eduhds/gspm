package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	Platform string
}

type GSConfig struct {
	GSPM     struct{}
	Packages []GSPackage
}

var args struct {
	Command string   `arg:"positional"`
	Values  []string `arg:"positional"`
}

const CONFIG_FILE = "config.json"

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

	configFile, err := os.ReadFile(CONFIG_FILE)

	if err != nil {
		tui.ShowError("Failed to read config file: " + err.Error())
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			tui.ShowError("Failed to parse config file: " + err.Error())
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
			tui.ShowInfo("No packages found")
			return
		}
		tui.ShowInfo(fmt.Sprintf("Loaded %d packages", countPackages))
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

			if gsPackage.Tag != "latest" {
				for _, configPackage := range config.Packages {
					if configPackage.Platform != runtime.GOOS {
						continue
					}
					if configPackage.Name == gsPackage.Name {
						if gsPackage.Tag == "" || configPackage.Tag == gsPackage.Tag {
							gsPackage = configPackage
						}
						break
					}
				}
			}

			var assetName = ""

			if gsPackage.AssetUrl == "" {
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
						if gsPackage.Tag == "latest" || gsPackage.Tag == release.TagName {
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

				assetName = tui.ShowOptions("Select an asset", assetOptions)

				for _, asset := range assets[tag] {
					if asset.Name == assetName {
						gsPackage.AssetUrl = asset.BrowserDownloadUrl
						break
					}
				}
			} else {
				downloadUrlSplited := strings.Split(gsPackage.AssetUrl, "/")
				assetNameFromUrl := downloadUrlSplited[len(downloadUrlSplited)-1]
				assetName = assetNameFromUrl
			}

			stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

			success := gitservice.GetGitHubReleaseAsset(assetName, gsPackage.AssetUrl)

			if success {
				stopAssetSpn("success")

				if gsPackage.Script == "" {
					runScript := tui.ShowConfirm("Do you want to run a script?")

					if runScript {
						tui.ShowInfo("Tip: Use $ASSET_PATH to reference the asset path")
						script := tui.ShowTextInput("Enter a script")
						gsPackage.Script = script

						if RunScript("~/Downloads/"+assetName, gsPackage.Script) {
							gsPackage.Platform = runtime.GOOS
							config.Packages = append(config.Packages, gsPackage)
						}
					} else {
						tui.ShowInfo("Script not provided.")
						tui.ShowSuccess("Package located at ~/Downloads/" + assetName)
					}
				} else {
					RunScript("~/Downloads/"+assetName, gsPackage.Script)
				}
			} else {
				stopAssetSpn("fail")
			}
		}

		configBytes, _ := json.MarshalIndent(config, "", "    ")

		err := os.WriteFile(CONFIG_FILE, configBytes, 0644)

		if err != nil {
			tui.ShowError(err.Error())
		}
	} else {
		tui.ShowError("Unknown command: " + args.Command)
		os.Exit(1)
	}
}
