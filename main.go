package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/eduhds/gspm/internal/gitservice"
	"github.com/eduhds/gspm/internal/tui"
	"github.com/eduhds/gspm/internal/util"
)

const appname = "gspm"
const version = "0.1.1"
const description = "Thanks for using gspm, the Git Services Package Manager.\n"
const asciiArt = "\n ,adPPYb,d8  ,adPPYba,  8b,dPPYba,   88,dPYba,,adPYba,  \n" +
	"a8\"    `Y88  I8[    \"\"  88P'    \"8a  88P'   \"88\"    \"8a \n" +
	"8b       88   `\"Y8ba,   88       d8  88      88      88 \n" +
	"\"8a,   ,d88  aa    ]8I  88b,   ,a8\"  88      88      88 \n" +
	" `\"YbbdP\"Y8  `\"YbbdP\"'  88`YbbdP\"'   88      88      88 \n" +
	" aa,    ,88             88                              \n" +
	"  \"Y8bbdP\"              88                              \n"

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

type args struct {
	Command string   `arg:"positional" help:"Command to run. Must be add, remove, update, install, edit, or list."`
	Repos   []string `arg:"positional" help:"Repos from Git Services (GitHub supported only for now). Format: username/repository"`
	Scripts []string `arg:"-s,--script,separate" help:"Script to run after download a asset. Use {{ASSET}} to reference the asset path."`
}

func (args) Version() string {
	return fmt.Sprintf("%s v%s", appname, version)
}

func (args) Description() string {
	return description
}

func (args) Epilogue() string {
	return "For more information visit https://github.com/eduhds/gspm"
}

var downloadPrefix = GetDownloadsDir()

func GetDownloadsDir() string {
	directory := filepath.Join(util.GetHomeDir(), "Downloads")
	util.CreateDirIfNotExists(directory)
	return directory
}

func GetConfigDir() string {
	directory := filepath.Join(util.GetHomeDir(), ".config")
	util.CreateDirIfNotExists(directory)
	return directory
}

func RunScript(assetName string, providedScript string) bool {
	assetPath := filepath.Join(downloadPrefix, assetName)

	script := strings.ReplaceAll(providedScript, "{{ASSET}}", assetPath)

	tui.ShowBox(fmt.Sprintf("ASSET=%s\n%s", assetPath, providedScript))

	execCommand, execOption := util.GetShellExec()

	out, err := exec.Command(execCommand, execOption, script).Output()

	if err != nil {
		if string(out) != "" {
			tui.ShowBox(string(out))
		}
		tui.ShowError(err.Error())
		return false
	} else {
		tui.ShowSuccess("Script executed successfully.")
		return true
	}
}

func ReadConfig() GSConfig {
	var config GSConfig

	configFile, err := os.ReadFile(filepath.Join(GetConfigDir(), fmt.Sprintf("%s.json", appname)))

	if err != nil {
		tui.ShowWarning("Config file not found: " + err.Error())
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			tui.ShowWarning("Config file is not valid: " + err.Error())
		}
	}

	os.WriteFile(filepath.Join(GetConfigDir(), fmt.Sprintf("%s.json.bkp", appname)), configFile, 0644)

	return config
}

func WriteConfig(config GSConfig) {
	configBytes, _ := json.MarshalIndent(config, "", "    ")

	err := os.WriteFile(filepath.Join(GetConfigDir(), fmt.Sprintf("%s.json", appname)), configBytes, 0644)

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
	var args args
	arg.MustParse(&args)

	if args.Command == "" {
		// Interactive mode
		tui.TextInfo(asciiArt)
		tui.ShowInfo(fmt.Sprintf("v%s", version))
		tui.ShowLine()

		quit := false

		for !quit {
			option := tui.ShowOptions("What command do you want to use?", []string{"add", "remove", "update", "install", "edit", "list"})

			repo := ""

			if option != "list" && option != "install" {
				repo = tui.ShowTextInput(fmt.Sprintf("What repository do you want \"%s\"? (Format: username/repository)", option), false, "")
			}

			program := strings.TrimSpace(strings.Join(os.Args, " "))

			cmd := exec.Command(program, option, repo)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()

			if err != nil {
				tui.ShowError(err.Error())
			}

			tui.ShowLine()
			quit = !tui.ShowConfirm("Continue to another command?")

			if quit {
				tui.TextSuccess(description)
				tui.ShowLine()
			}
		}

		return
	}

	// Non-interactive mode
	tui.ShowInfo(fmt.Sprintf("%s v%s", appname, version))
	tui.ShowLine()

	config := ReadConfig()
	platformPackages := PlatformPackages(config)
	countPackages := len(platformPackages)

	if args.Command == "install" {
		if countPackages == 0 {
			tui.ShowInfo("No packages found")
			return
		}

		if len(args.Repos) > 0 || len(args.Scripts) > 0 {
			tui.ShowWarning("Ignoring args for install command.")
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
		if len(args.Repos) == 0 {
			tui.ShowError("No packages provided")
			return
		}

		for index, value := range args.Repos {
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

				if len(args.Scripts) > 0 && args.Scripts[index] != "" {
					gsPackage.Script = args.Scripts[index]
				}

				if gsPackage.Script == "" {
					runScript := tui.ShowConfirm("Do you want to run a script?")

					if runScript {
						tui.ShowInfo("Use {{ASSET}} to reference the asset path")
						gsPackage.Script = tui.ShowTextInput("Enter a script", true, "")

						if RunScript(assetName, gsPackage.Script) {
							gsPackage.Platform = runtime.GOOS
							if args.Command == "update" {
								for index, configPackage := range config.Packages {
									if configPackage.Platform == runtime.GOOS && configPackage.Name == gsPackage.Name {
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
						tui.ShowSuccess("Package located at " + filepath.Join(downloadPrefix, assetName))
					}
				} else {
					if RunScript(assetName, gsPackage.Script) {
						if args.Command == "update" {
							for index, configPackage := range config.Packages {
								if configPackage.Platform == runtime.GOOS && configPackage.Name == gsPackage.Name {
									config.Packages[index] = gsPackage
									break
								}
							}
						} else {
							config.Packages = append(config.Packages, gsPackage)
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
			if item.Platform == runtime.GOOS && item.Name == args.Repos[0] {
				tui.ShowInfo("Package " + item.Name + " removed")
			} else {
				keepedPackages = append(keepedPackages, item)
			}
		}

		if len(keepedPackages) == len(config.Packages) {
			tui.ShowError("Package not found: " + args.Repos[0])
			os.Exit(1)
		}

		config.Packages = keepedPackages

		WriteConfig(config)
	} else if args.Command == "edit" {
		if countPackages == 0 {
			tui.ShowInfo("No packages found")
			return
		}

		for index, value := range args.Repos {
			if index > 0 {
				tui.ShowLine()
			}
			packageName := value

			var gsPackage GSPackage

			for _, configPackage := range platformPackages {
				if configPackage.Name == packageName {
					gsPackage = configPackage
					break
				}
			}

			if gsPackage.Name == "" {
				tui.ShowError("Package not found: " + packageName)
				continue
			}

			tui.ShowWarning("Editing script for package " + gsPackage.Name)

			tui.ShowInfo("Use {{ASSET}} to reference the asset path")
			script := tui.ShowTextInput("Enter a script", true, gsPackage.Script)
			gsPackage.Script = script

			for index, configPackage := range config.Packages {
				if configPackage.Platform == runtime.GOOS && configPackage.Name == gsPackage.Name {
					config.Packages[index] = gsPackage
					break
				}
			}
		}

		WriteConfig(config)
	} else {
		tui.ShowError("Unknown command: " + args.Command)
		os.Exit(1)
	}
}
