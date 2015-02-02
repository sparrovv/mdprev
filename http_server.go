package mdprev

import (
	"fmt"
	"log"
	"net/http"
)

// start HTTP server
func (mdPrev *MdPrev) RunServer(portNumber string) {
	http.Handle("/"+mdPrev.MdFile, mdFileHandler(mdPrev))
	http.Handle("/", staticFileHandler(mdPrev))

	hub := newHub(mdPrev.Broadcast, mdPrev.Exit)
	http.Handle("/ws", wsHandler(hub))
	go hub.run()

	log.Fatal(http.ListenAndServe(":"+portNumber, nil))
}

func mdFileHandler(mdPrev *MdPrev) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := ToHTML(mdPrev.MdContent)

		fmt.Fprintf(w, html.String())
	})
}

// handle static files, e.g. images
func staticFileHandler(mdPrev *MdPrev) http.Handler {
	return http.FileServer(http.Dir(mdPrev.MdDirPath()))
}

// WebSocket handler. Register all clients to the hub, that sends updates if file changed
func wsHandler(h *hub) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := NewWSConnection(ws)
		h.register <- c
		go c.writer()
		c.unregisterOnEOF(h)
	})
}
