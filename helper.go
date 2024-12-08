package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/eduhds/gspm/internal/tui"
	"github.com/eduhds/gspm/internal/util"
)

func PlatformPackages(config GSConfig) []GSPackage {
	var packages []GSPackage
	for _, item := range config.Packages {
		if item.Platform == runtime.GOOS {
			packages = append(packages, item)
		}
	}
	return packages
}

func ResolvePackage(value string, cfg GSConfig, mustExist bool) (GSPackage, error) {
	packages := PlatformPackages(cfg)

	info := strings.Split(value, "@")
	name := info[0]
	tag := ""
	platform := runtime.GOOS

	if len(strings.Split(name, "/")) != 2 {
		return GSPackage{}, errors.New("Invalid package name: " + name)
	}

	if len(info) > 1 {
		tag = info[1]
	}

	for _, item := range packages {
		if item.Name == name {
			if tag != "" && item.Tag != tag {
				item.Tag = tag
			}
			return item, nil
		}
	}

	if mustExist {
		return GSPackage{}, errors.New("Package not found: " + name)
	}

	return GSPackage{Name: name, Tag: tag, Platform: platform}, nil
}

func RunScript(assetName string, providedScript string) bool {
	assetPath := filepath.Join(downloadPrefix, assetName)

	script := strings.ReplaceAll(providedScript, "{{ASSET}}", assetPath)

	tui.ShowBox(fmt.Sprintf("ASSET=%s\n%s", assetPath, providedScript))

	execCommand, execOption := util.GetShellExec()

	cmd := exec.Command(execCommand, execOption, script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		tui.ShowError(err.Error())
		return false
	} else {
		tui.ShowSuccess("Script executed successfully.")
		return true
	}
}

func ResolveConfig() GSConfig {
	var config GSConfig

	configFile, err := os.ReadFile(filepath.Join(util.GetConfigDir(customConfigDir), fmt.Sprintf("%s.json", appname)))

	if err != nil {
		tui.ShowWarning("Config file not found: " + err.Error())
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			tui.ShowWarning("Config file is not valid: " + err.Error())
		}
	}

	os.WriteFile(filepath.Join(util.GetConfigDir(customConfigDir), fmt.Sprintf("%s.json.bkp", appname)), configFile, 0644)

	return config
}

func WriteConfig(config GSConfig) {
	config.GSPM.Version = "1"

	configBytes, _ := json.MarshalIndent(config, "", "    ")

	err := os.WriteFile(filepath.Join(util.GetConfigDir(customConfigDir), fmt.Sprintf("%s.json", appname)), configBytes, 0644)

	if err != nil {
		tui.ShowWarning("Cannot to write config file: " + err.Error())
	}
}

func AssetNameFromUrl(url string) string {
	urlSplited := strings.Split(url, "/")
	assetNameFromUrl := urlSplited[len(urlSplited)-1]
	return assetNameFromUrl
}
