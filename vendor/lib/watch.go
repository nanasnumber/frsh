package lib

import (
	"fmt"
	"os"
	"time"
)

func Watch(p string, fn func()) error {

	path := p

	// fmt.Println(path)

	initStat, err := os.Stat(path)

	if err != nil {
		// fmt.Println(err)
	}

	for {
		stat, err := os.Stat(path)
		if err != nil {
			// fmt.Println("file not found")
		}

		if stat != nil {
			if stat.Size() != initStat.Size() || stat.ModTime() != initStat.ModTime() {
				fmt.Println("===============")
				fmt.Println(path + " changed")
				fn()
				Watch(path, fn)
			}

			time.Sleep(500 * time.Millisecond)
		}

	}

	return nil
}
