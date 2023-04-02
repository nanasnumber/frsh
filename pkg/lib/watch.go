/*
* File watcher with callback pattern
* */

package lib

import (
	"fmt"
	"frsh/pkg/io"
	"os"
	"time"
)

func Watch(p string, fn func()) error {
	path := p
	initStat, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}

	initCheckSum := io.CheckSum(path)

	for {
		curStat, err := os.Stat(path)
		if err != nil {
			fmt.Println(err)
		}

		curCheckSum := io.CheckSum(path)

		if curCheckSum != initCheckSum || curStat.ModTime() != initStat.ModTime() {
			fmt.Println("===============")
			fmt.Println(path + " changed")
			fn()
			Watch(path, fn)
		}

		time.Sleep(500 * time.Millisecond)

	}

	return nil
}
