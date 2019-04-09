package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	clientIdentities map[string]*websocket.Conn
)

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userName := r.Header.Get("X-Small-Chat-Id")
	clientIdentities[userName] = conn

	// Read messages from socket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			return
		}

		log.Printf("msg from %s: %s", string(userName), string(msg))

		switch userName {
		case "vic":
			go clientIdentities["judy"].WriteMessage(websocket.TextMessage, msg)
		case "judy":
			go clientIdentities["vic"].WriteMessage(websocket.TextMessage, msg)
		}

	}
}

func newUserMap() {
	clientIdentities = map[string]*websocket.Conn{}
}

func main() {
	newUserMap()
	http.HandleFunc("/", ws)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
