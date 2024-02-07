package gitservice

import (
	"errors"
	"time"

	"github.com/imroc/req/v3"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

var client = req.C().
	SetTimeout(5 * time.Second)

func GetGitHubReleases(username string, repository string) ([]GSGitHubRelease, error) {
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
