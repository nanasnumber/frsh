package lib

import "path/filepath"

func FileIgnore(p string) bool {

	ext := filepath.Ext(p)

	switch ext {
	case "":
		return true
	case ".swp":
		return true
	default:
		return false
	}
}
