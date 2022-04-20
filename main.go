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

func wcClient() string {
	tpl := `
		<script>
			(() => {

				const livereload = () => {
					const c = new WebSocket("ws://localhost:%s/livereload");

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
	port := setPort()
	script := fmt.Sprintf(tpl, port)
	return script
}

func fc(p string) string {
	file, err := os.ReadFile(p)

	if err != nil {
		fmt.Println(p + " not exist")
	}

	fileContent := string(file)

	return fileContent
}

func fileResponse(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	m := lib.MIMEType
	mimeType := m(filePath)

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
