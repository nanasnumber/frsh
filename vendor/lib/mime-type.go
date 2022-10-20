/*
* Return MIME type based on file extension
* */

package lib

import (
	"path/filepath"
)

func MIMEType(p string) string {

	ext := filepath.Ext(p)

	mimeMap := make(map[string]string)

	mimeMap[""] = "text/html"
	mimeMap[".html"] = "text/html"
	mimeMap[".css"] = "text/css"
	mimeMap[".js"] = "text/javascript"
	mimeMap[".json"] = "application/json"
	mimeMap[".svg"] = "image/svg+xml"
	mimeMap[".jpg"] = "image/jpeg"
	mimeMap[".ico"] = "image/x-icon"

	return mimeMap[ext]
}
