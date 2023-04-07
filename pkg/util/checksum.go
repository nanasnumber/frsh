package util

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

/* Get file checksum by file path */

func CheckSum(path string) string {
	file, err := os.ReadFile(path)

	if err != nil {
		fmt.Println(err)
	}

	hash := sha1.New()
	hash.Write([]byte(file))
	response := hex.EncodeToString(hash.Sum(nil))

	return response
}
