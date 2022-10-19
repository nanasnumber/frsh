/*
* Websockets client
* */
package lib

import "fmt"

func WebsocketsClient(p string) string {
	tpl := `
		<script>
			(() => {

				const livereload = () => {
					const c = new WebSocket("ws://localhost:%s/livereload");

					c.onopen =  (e) => {
						console.info('FRSH: Livereload Connected');
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
							console.info('FRSH: Livereload Disconnected');
							c.close()
						}

						c.onopen =  () => {
							console.info('FRSH: Livereload Connected');
							c.send('pong');
						};
					}
				};

				livereload();
			})();
		</script>
	`
	port := p
	script := fmt.Sprintf(tpl, port)
	return script
}
