package main

import (
	"fmt"
	"frsh/pkg/lib"
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

		//NOTE:
		//just a quick patch to fix for recognize *.html as directory
		//instead of file
		ext := filepath.Ext(r.URL.Path)
		if ext == ".html" {
			filePath = "." + r.URL.Path
		} else {
			filePath = "." + r.URL.Path + "/index.html"
		}

		contentOrigin := fc(filePath)
		content := strings.Replace(contentOrigin, "</head>", wc(p)+"</head>", -1)
		w.Header().Set("Content-Type", mimeType)
		w.Write([]byte(content))
	} else if fc(filePath) == "404" {
		http.NotFound(w, r)
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

	var wg sync.WaitGroup

	http.HandleFunc("/livereload", func(w http.ResponseWriter, r *http.Request) {

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			// fmt.Println(err)
		}

		fileWalkCallback := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ignoreDir := lib.DirIgnore
			if ignoreDir(path) {
				return nil
			}

			ignoreFile := lib.FileIgnore
			if ignoreFile(path) {
				return nil
			}

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

			return nil
		}

		fileErr := filepath.Walk("./", fileWalkCallback)

		wg.Wait()

		if fileErr != nil {
			fmt.Println(fileErr)
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
