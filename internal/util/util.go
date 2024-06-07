package util

import (
	"os"
	"os/user"
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
