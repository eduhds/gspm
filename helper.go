package main

import (
	"errors"
	"strings"
)

func ResolvePackage(repo string, packages []GSPackage) (GSPackage, error) {
	for _, item := range packages {
		if item.Name == repo {
			return item, nil
		}
	}

	info := strings.Split(repo, "@")
	name := info[0]

	if len(strings.Split(name, "/")) != 2 {
		return GSPackage{}, errors.New("Invalid package name: " + name)
	}

	tag := ""

	if len(info) > 1 {
		tag = info[1]
	}

	return GSPackage{Name: name, Tag: tag}, nil
}
