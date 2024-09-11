package util

import (
	"os"
	"os/user"
	"runtime"
)

func GetHomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
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
