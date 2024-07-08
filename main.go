package main

import (
	"flag"

	"github.com/gorilla/websocket"
)

type MysteryGiftDB struct {
}

var upgrader = websocket.Upgrader{}
var addr = flag.String("addr", ":8080", "http service address")
var mgdb MysteryGiftDB

func main() {
	flag.Parse()
	mgdb.SetUpDB()
	mgdb.InsertDBData()
}
