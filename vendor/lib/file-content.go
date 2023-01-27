/*
* Get file content by path
* */
package lib

import (
	"fmt"
	"os"
)

func FileContent(p string) string {
	file, err := os.ReadFile(p)

	if err != nil {
		fmt.Println(p + " does not exist")
		fmt.Println(err)
		return "404"
	}

	fileContent := string(file)

	return fileContent
}
