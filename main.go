package main

import (
	"flag"
	"fmt"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var addr = flag.String("addr", ":8080", "http service address")

func main() {
	fmt.Print("Set up db!\n")
	SetUpDB()
	fmt.Print("Have successfully set up db!\nInserting data now!\n")
	InsertDBData()
	fmt.Print("Have successfully inserted data!\n")
}
