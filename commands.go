package main

import (
  "fmt"
  "runtime"

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
  fmt.Println("Name: " + gsp.Name)
  fmt.Println("Tag: " + gsp.Tag)
  fmt.Println("AssetUrl: " + gsp.AssetUrl)
  fmt.Println("Script: " + gsp.Script)
  fmt.Println("Platform: " + gsp.Platform)
  fmt.Println("LastModfied: " + gsp.LastModfied)
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

