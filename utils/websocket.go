package utils

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	Reader(ws)
}

func HandleRoutes() {
	var addr = flag.String("addr", ":8080", "http service address")
	http.HandleFunc("/ws", WsEndpoint)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("Error occured while setting up server!\nError is %v", err)
	}
}
