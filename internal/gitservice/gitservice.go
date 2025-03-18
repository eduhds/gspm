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
	SetTimeout(10 * 60 * time.Second)

// Deprecated: Use sepcific service function
func GetGitHubReleases(packageName string) ([]GSGitHubRelease, error) {
	var username string = strings.Split(packageName, "/")[0]
	var repository string = strings.Split(packageName, "/")[1]
	var releases []GSGitHubRelease
	var errMsg ErrorMessage

	cr := client.R()

	if GHToken != "" {
		cr.SetHeader("Authorization", fmt.Sprintf("Bearer %s", GHToken))
	}

	resp, err := cr.
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

// Deprecated: Use sepcific service function
func GetGitHubReleaseAsset(assetName string, assetDownloadUrl string) (bool, error) {
	outputFile := filepath.Join(util.GetHomeDir(), "Downloads", assetName)
	var errMsg ErrorMessage

	cr := client.R()

	if GHToken != "" {
		cr.SetHeader("Authorization", fmt.Sprintf("Bearer %s", GHToken))
	}

	resp, err := cr.SetOutputFile(outputFile).
		SetHeader("Accept", "application/octet-stream").
		SetErrorResult(&errMsg).
		Get(assetDownloadUrl)

	if err != nil {
		return false, err
	}

	if resp.IsErrorState() {
		return false, errors.New(errMsg.Message)
	}

	return true, nil
}

func GSGetReleases(url string, token string, releases any) (bool, error) {
	var errMsg ErrorMessage

	cr := client.R()

	if token != "" {
		cr.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := cr.
		SetSuccessResult(&releases).
		SetErrorResult(&errMsg).
		EnableDump().
		Get(url)

	if err != nil { // Error handling.
		return false, err
	}

	if resp.IsErrorState() { // Status code >= 400.
		return false, errors.New(errMsg.Message)
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		return true, nil
	}

	// Unknown status code.
	return false, errors.New("Unknown status: " + resp.Status)
}

func GSGetReleaseAsset(assetName string, assetDownloadUrl string, token string) (bool, error) {
	outputFile := filepath.Join(util.GetHomeDir(), "Downloads", assetName)
	var errMsg ErrorMessage

	cr := client.R()

	if token != "" {
		cr.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := cr.SetOutputFile(outputFile).
		SetHeader("Accept", "application/octet-stream").
		SetErrorResult(&errMsg).
		Get(assetDownloadUrl)

	if err != nil {
		return false, err
	}

	if resp.IsErrorState() {
		return false, errors.New(errMsg.Message)
	}

	return true, nil
}
