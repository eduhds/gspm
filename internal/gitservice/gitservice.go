package gitservice

import (
	"fmt"
	"log"
	"time"

	"github.com/imroc/req/v3"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

var client = req.C().
	SetTimeout(5 * time.Second)

func GetGitHubReleases(username string, repository string) {
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
		log.Println("error:", err)
		log.Println("raw content:")
		log.Println(resp.Dump()) // Record raw content when error occurs.
		return
	}

	if resp.IsErrorState() { // Status code >= 400.
		fmt.Println(errMsg.Message) // Record error message returned.
		return
	}

	if resp.IsSuccessState() { // Status code is between 200 and 299.
		fmt.Printf("%s (%s)\n", releases[0].Name, releases[0].Body)
		return
	}

	// Unknown status code.
	log.Println("unknown status", resp.Status)
	log.Println("raw content:")
	log.Println(resp.Dump()) // Record raw content when server returned unknown status code.
}
