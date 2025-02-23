package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

type TimeResponse struct {
	Time string `json:"time"`
}

func Echo(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{} // use default options
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var timeResponse TimeResponse
		if err := json.Unmarshal(message, &timeResponse); err != nil {
			log.Println("unmarshal:", err)
			break
		}

		log.Printf("recv: %s", timeResponse.Time)
		gift := SearchDBData(timeResponse.Time)

		// Convert gifts to JSON bytes
		giftBytes, err := json.Marshal(gift)
		if err != nil {
			log.Println("marshal:", err)
			break
		}

		err = c.WriteMessage(mt, giftBytes)
		fmt.Println("Sent all active events!")
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func HandleRoutes() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/mysterygift", Echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
	fmt.Println("Now running!")
}
