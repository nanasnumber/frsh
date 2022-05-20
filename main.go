package main

import (
	"fmt"
	"lib"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

func setPort() string {
	return "8080"
}

func fileResponse(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	m := lib.MIMEType
	fc := lib.FileContent
	wc := lib.WebsocketsClient
	p := setPort()
	mimeType := m(filePath)

	if mimeType == "text/html" {
		filePath := "." + r.URL.Path + "/index.html"
		contentOrigin := fc(filePath)
		content := strings.Replace(contentOrigin, "</head>", wc(p)+"</head>", -1)
		w.Header().Set("Content-Type", mimeType)
		w.Write([]byte(content))
	} else {
		filePath := "." + r.URL.Path
		content := fc(filePath)
		w.Header().Set("Content-Type", mimeType)
		w.Write([]byte(content))
	}
}

func watchAndReload() {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/livereload", func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// fmt.Println(err)
		}

		var wg sync.WaitGroup

		fileErr := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)

			/*
			* exclude temp file and directory for watcher
			* [TODO]: Need a method to check file extentions to ignore reload
			* */
			if ext != ".swp" && ext != "" {

				wg.Add(1)

				go func() {

					defer wg.Done()

					w := lib.Watch
					w(path, func() {

						for {
							msgType, _, err := c.ReadMessage()
							if err != nil {
								// fmt.Println(err)
							}

							err = c.WriteMessage(msgType, []byte("ping"))
							if err != nil {
								return
							}
						}
					})

				}()
			}

			return nil
		})

		wg.Wait()

		if fileErr != nil {
			// fmt.Println(fileErr)
		}
	})
}

func main() {
	watchAndReload()
	port := setPort()
	http.HandleFunc("/", fileResponse)
	fmt.Println("Listening to port:" + port)
	http.ListenAndServe(":"+port, nil)
}
