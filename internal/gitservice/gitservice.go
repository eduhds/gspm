package gitservice

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/eduhds/gspm/internal/util"
	"github.com/imroc/req/v3"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

var client = req.C().
	SetTimeout(60 * time.Second)

func GetGitHubReleases(packageName string) ([]GSGitHubRelease, error) {
	var username string = strings.Split(packageName, "/")[0]
	var repository string = strings.Split(packageName, "/")[1]
	var releases []GSGitHubRelease
	var errMsg ErrorMessage

	resp, err := client.R().
		SetPathParam("username", username).
		SetPathParam("repo", repository).
		SetSuccessResult(&releases).
		SetErrorResult(&errMsg).
		EnableDump().
		Get("https://api.github.com/repos/{username}/{repo}/releases")

	if err != nil { // Error handling.
		return nil, err
	}

	if resp.IsErrorState() { // Status code >= 400.
		return nil, errors.New(errMsg.Message)
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		return releases, nil
	}

	// Unknown status code.
	return nil, errors.New("Unknown status: " + resp.Status)
}

func GetGitHubReleaseAsset(assetName string, assetDownloadUrl string) bool {
	outputFile := filepath.Join(util.GetHomeDir(), "Downloads", assetName)

	_, err := client.R().SetOutputFile(outputFile).
		Get(assetDownloadUrl)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
