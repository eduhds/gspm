package main

import (
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

var args struct {
	Command string   `arg:"positional"`
	Values  []string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)

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

			for _, asset := range assets[tag] {
				assetOptions = append(assetOptions, asset.Name)
			}

			assetName := tui.ShowOptions("Select an asset", assetOptions)

			stopAssetSpn := tui.ShowSpinner(fmt.Sprintf("Downloading %s...", assetName))

			var res bool

			for _, asset := range assets[tag] {
				if asset.Name == assetName {
					gsPackage.AssetUrl = asset.BrowserDownloadUrl
					res = gitservice.GetGitHubReleaseAsset(asset)
					break
				}
			}

			if res {
				stopAssetSpn("success")
				tui.ShowInfo("Asset downloaded path: /tmp/" + assetName)

				runScript := tui.ShowConfirm("Do you want to run a script?")

				if runScript {
					script := tui.ShowTextInput("Enter a script")
					gsPackage.Script = script
				}

				if len(gsPackage.Script) == 0 {
					tui.ShowInfo("Script not provided.")

					_, err := exec.Command("bash", "-c", "cp /tmp/"+assetName+" ~/Downloads").Output()

					if err != nil {
						tui.ShowError(err.Error())
					} else {
						tui.ShowSuccess("Asset copied to ~/Downloads/" + assetName)
					}
				} else {
					stopScriptSpn := tui.ShowSpinner("Running provided script...")

					tui.ShowBox(gsPackage.Script)

					out, err := exec.Command("bash", "-c", gsPackage.Script).Output()

					if err != nil {
						stopScriptSpn("fail")
						tui.ShowBox(string(out))
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
	} else {
		tui.ShowError("Unknown command: " + args.Command)
		os.Exit(1)
	}
}
