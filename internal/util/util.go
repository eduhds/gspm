package util

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
)

func GetHomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}

func GetConfigDir(overrideDir string) string {
	directory := filepath.Join(GetHomeDir(), ".config")

	if overrideDir != "" {
		directory = overrideDir
	}

	CreateDirIfNotExists(directory)
	return directory
}

func GetDownloadsDir() string {
	directory := filepath.Join(GetHomeDir(), "Downloads")
	CreateDirIfNotExists(directory)
	return directory
}

func CreateDirIfNotExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

func GetShellExec() (string, string) {
	execCommand := "/bin/sh"
	execOption := "-c"

	if runtime.GOOS == "windows" {
		execCommand = "cmd.exe"
		execOption = "/C"
	}

	return execCommand, execOption
}

func DateLocalNow() string {
	return time.Now().Local().Format(time.RFC3339)
}
