package lib

import (
	"strings"
)

func DirIgnore(p string) bool {

	var list = []string{
		".git",
		"node_modules",
	}

	for _, item := range list {
		if strings.Contains(p, item) {
			return true
		}
	}

	return false
}
