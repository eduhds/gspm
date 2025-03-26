package gitservice

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/eduhds/gspm/internal/util"
	"github.com/imroc/req/v3"
)

const (
	GITHUB    string = "github"
	GITLAB    string = "gitlab"
	BITBUCKET string = "bitbucket"
)

var ServiceSymbol = map[string]string{
	"gitlab":    "🦊",
	"github":    "🐙",
	"bitbucket": "🪣 ",
}

var client = req.C().
	SetTimeout(10 * 60 * time.Second)

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
		//SetHeader("Accept", "application/octet-stream").
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
