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
		fmt.Println(p + " not exist")
	}

	fileContent := string(file)

	return fileContent
}
