package util

import "os/user"

func GetHomeDir() string {
	u, _ := user.Current()
	return u.HomeDir
}
