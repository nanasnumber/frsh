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

func fc(p string) string {

	file, err := os.ReadFile(p)

	if err != nil {
		fmt.Println(err)
	}

	fileContent := string(file)

	return fileContent
}

func wcClient() string {
	return `
		<script>
			(() => {

				const livereload = () => {
					const c = new WebSocket("ws://localhost:8080/livereload");

					c.onopen =  (e) => {
						console.info('Livereload Connected');
						c.send('pong');
					};

					c.onmessage = (e) => {
						if (e.data === "ping") {
							location.reload();
						}
					};

					c.onclose = (e) => {

						console.log(e)

						if (e.type === 'close') {
							c.close()
						}

						c.onopen =  () => {
							console.info('Livereload Connected');
							c.send('pong');
						};
					}
				};

				livereload();
			})();
		</script>
	`
}

func MIMEType(p string) string {
	ext := filepath.Ext(p)

	switch ext {
	case "":
		return "text/html"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
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

func fileResponse(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	mimeType := MIMEType(filePath)

	if mimeType == "text/html" {
		filePath := "." + r.URL.Path + "/index.html"
		contentOrigin := fc(filePath)
		content := strings.Replace(contentOrigin, "</head>", wcClient()+"</head>", -1)
		w.Header().Set("Content-Type", mimeType)
		w.Write([]byte(content))
	} else {
		filePath := "." + r.URL.Path
		content := fc(filePath)
		w.Header().Set("Content-Type", mimeType)
		w.Write([]byte(content))
	}

}

func main() {

	go func() {

		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		http.HandleFunc("/livereload", func(w http.ResponseWriter, r *http.Request) {

			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err)
			}

			/*
			* TODO: not fully sure how wait group supposed to work yet, need to study this one
			* */
			var wg sync.WaitGroup

			fileErr := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				ext := filepath.Ext(path)

				/*
				* exclude temp file and directory for watcher
				* */
				if ext != ".swp" && ext != "" {

					wg.Add(1)

					go func() {
						defer wg.Done()
						w := lib.Watch
						w(path, func() {

							if err != nil {
								fmt.Println(err)
							}

							for {

								msgType, msg, err := c.ReadMessage()
								if err != nil {
									fmt.Println(err)
								}

								fmt.Println(string(msg))

								if string(msg) == "pong" {
									fmt.Println("client responded")
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
				fmt.Println(fileErr)
			}
		})

	}()

	go func() {
		http.HandleFunc("/", fileResponse)
		http.ListenAndServe(":8080", nil)
	}()

	select {}

}
