package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
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

var downloadPrefix = GetDownloadsDir()

func CreateDirIfNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

func GetHomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}

func GetDownloadsDir() string {
	directory := GetHomeDir() + "/Downloads"
	CreateDirIfNotExists(directory)
	return directory
}

func GetConfigDir() string {
	directory := GetHomeDir() + "/.config"
	CreateDirIfNotExists(directory)
	return directory
}

func RunScript(assetName string, providedScript string) bool {
	assetPath := downloadPrefix + "/" + assetName
	script := "ASSET=" + assetPath + "\n" + providedScript

	stopScriptSpn := tui.ShowSpinner("Running provided script...")

	tui.ShowBox(script)

	out, err := exec.Command("/bin/sh", "-c", script).Output()

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

	configFile, err := os.ReadFile(GetConfigDir() + "/gspm.json")

	if err != nil {
		tui.ShowWarning("Config file not found: " + err.Error())
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			tui.ShowWarning("Config file is not valid: " + err.Error())
		}
	}

	return config
}

func WriteConfig(config GSConfig) {
	configBytes, _ := json.MarshalIndent(config, "", "    ")

	err := os.WriteFile(GetConfigDir()+"/gspm.json", configBytes, 0644)

	if err != nil {
		tui.ShowWarning("Cannot to write config file: " + err.Error())
	}
}

func PlatformPackages(config GSConfig) []GSPackage {
	var packages []GSPackage
	for _, item := range config.Packages {
		if item.Platform == runtime.GOOS {
			packages = append(packages, item)
		}
	}
	return packages
}

func AssetNameFromUrl(url string) string {
	urlSplited := strings.Split(url, "/")
	assetNameFromUrl := urlSplited[len(urlSplited)-1]
	return assetNameFromUrl
}

func main() {
	arg.MustParse(&args)

	config := ReadConfig()
	platformPackages := PlatformPackages(config)
	countPackages := len(platformPackages)

	if args.Command == "install" {
		if countPackages == 0 {
			tui.ShowInfo("No packages found")
			return
		}
		tui.ShowInfo(fmt.Sprintf("Loaded %d packages", countPackages))

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
	} else if args.Command == "list" {
		if countPackages == 0 {
			tui.ShowInfo("No packages found")
			return
		}

		tui.ShowInfo(fmt.Sprintf("%d packages for %s", countPackages, runtime.GOOS))
		tui.ShowLine()

		for _, item := range platformPackages {
			tui.ShowMessage("📦 " + item.Name + "@" + item.Tag)
			tui.ShowMessage("🔗 " + item.AssetUrl)
			tui.ShowMessage("🛠️  " + item.Script)
			tui.ShowLine()
		}
	} else if args.Command == "add" || args.Command == "update" {
		if len(args.Values) == 0 {
			tui.ShowError("No packages provided")
			return
		}

		for index, value := range args.Values {
			if index > 0 {
				tui.ShowLine()
			}
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

			if args.Command == "update" {
				packageTag = "latest"
			}

			var gsPackage GSPackage
			gsPackage.Name = packageName
			gsPackage.Tag = packageTag

			var assetName = ""

			for _, configPackage := range platformPackages {
				if configPackage.Name == gsPackage.Name {
					if args.Command == "update" {
						assetName = AssetNameFromUrl(configPackage.AssetUrl)
					} else {
						if gsPackage.Tag != "latest" {
							if gsPackage.Tag == "" || configPackage.Tag == gsPackage.Tag {
								gsPackage.AssetUrl = configPackage.AssetUrl
							}
						}
					}
					gsPackage.Script = configPackage.Script
					gsPackage.Platform = configPackage.Platform
					break
				}
			}

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
					if assetName == asset.Name {
						gsPackage.AssetUrl = asset.BrowserDownloadUrl
						break
					}
					assetOptions = append(assetOptions, asset.Name)
				}

				if gsPackage.AssetUrl == "" {
					assetName = tui.ShowOptions("Select an asset", assetOptions)
					for _, asset := range assets[tag] {
						if asset.Name == assetName {
							gsPackage.AssetUrl = asset.BrowserDownloadUrl
							break
						}
					}
				}
			} else {
				assetName = AssetNameFromUrl(gsPackage.AssetUrl)
			}

			stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

			success := gitservice.GetGitHubReleaseAsset(assetName, gsPackage.AssetUrl)

			if success {
				stopAssetSpn("success")

				if gsPackage.Script == "" {
					runScript := tui.ShowConfirm("Do you want to run a script?")

					if runScript {
						tui.ShowInfo("Use $ASSET to reference the asset path")
						script := tui.ShowTextInput("Enter a script", true)
						gsPackage.Script = script

						if RunScript(assetName, gsPackage.Script) {
							gsPackage.Platform = runtime.GOOS
							if args.Command == "update" {
								for index, configPackage := range config.Packages {
									if configPackage.Name == gsPackage.Name {
										config.Packages[index] = gsPackage
										break
									}
								}
							} else {
								config.Packages = append(config.Packages, gsPackage)
							}
						}
					} else {
						tui.ShowInfo("Script not provided.")
						tui.ShowSuccess("Package located at " + downloadPrefix + "/" + assetName)
					}
				} else {
					if RunScript(assetName, gsPackage.Script) {
						if args.Command == "update" {
							for index, configPackage := range config.Packages {
								if configPackage.Name == gsPackage.Name {
									config.Packages[index] = gsPackage
									break
								}
							}
						}
					}
				}
			} else {
				stopAssetSpn("fail")
			}
		}

		WriteConfig(config)
	} else if args.Command == "remove" {
		if countPackages == 0 {
			tui.ShowInfo("No packages found")
			return
		}

		var keepedPackages []GSPackage

		for _, item := range config.Packages {
			if item.Platform == runtime.GOOS && item.Name == args.Values[0] {
				tui.ShowInfo("Package " + item.Name + " removed")
			} else {
				keepedPackages = append(keepedPackages, item)
			}
		}

		config.Packages = keepedPackages

		WriteConfig(config)
	} else {
		tui.ShowError("Unknown command: " + args.Command)
		os.Exit(1)
	}
}
