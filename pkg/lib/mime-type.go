/*
* Return MIME type based on file extension
* */

package lib

import "path/filepath"

func MIMEType(p string) string {
	ext := filepath.Ext(p)

	switch ext {
	case "":
		return "text/html"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".json":
		return "application/json"
	case ".js":
		return "text/javascript"
	case ".jpg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}
