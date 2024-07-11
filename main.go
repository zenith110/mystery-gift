package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/websocket"
	Models "github.com/zenith110/pokemon-go-engine-toml-models/models"
)

var upgrader = websocket.Upgrader{}
var addr = flag.String("addr", ":8080", "http service address")
var client *badger.DB

func SetUpDB() {
	if client != nil {
		var err error
		opt := badger.DefaultOptions("").WithInMemory(true)
		client, err = badger.Open(opt)
		if err != nil {
			fmt.Printf("Error while setting up db!\nError is %v", err)
		}
	}
}

func InsertDBData() {
	file, err := os.Open("gifts.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var mysteryGifts Models.MysteryGiftsServer

	err = toml.Unmarshal(bytes, &mysteryGifts)

	if err != nil {
		fmt.Printf("Error occured while reading data!\nError is %v", err)
		// panic(err)
	}
	var genericErr error
	for mysteryGiftIndex := range mysteryGifts.Mysterygifts {
		fmt.Printf("Looking at %s\n", mysteryGifts.Mysterygifts[mysteryGiftIndex].Name)
		err = client.Update(func(txn *badger.Txn) error {
			if mysteryGifts.Mysterygifts[mysteryGiftIndex].GiftType == "Pokemon" {
				fmt.Print("hello")
				// data := fmt.Sprintf("%v", mysteryGifts.Mysterygifts[mysteryGiftIndex].Pokemongift)
				// fmt.Print(data)
				// txn.Set([]byte(mysteryGifts.Mysterygifts[mysteryGiftIndex].Name), []byte(data))
				return nil
			}
			return genericErr
		})
		if err != nil {
			fmt.Printf("Error occured while storing data!\nError is %v", err)
			// panic(err)
		}
	}
}

func TimeCheck(startDate, endDate, playerDate time.Time) bool {
	return startDate.After(playerDate) && endDate.Before(playerDate)
}

func SearchDBData(currentDate string, giftName string) Models.MysteryGiftServer {
	// Create a byte array to copy the data from the db into
	var mysteryGiftData []byte

	// Create the struct to dump the returned data into
	var mysteryGift Models.MysteryGiftServer
	// Begins the transaction of viewing the data
	err := client.View(func(txn *badger.Txn) error {
		// Grabs the specific data based off the gift name
		data, _ := txn.Get([]byte(giftName))
		// Creates a copy to be used for use outside of the db view
		err := data.Value(func(val []byte) error {
			// Copying or parsing val is valid.
			mysteryGiftData = append([]byte{}, val...)

			return nil
		})

		if err != nil {
			panic(err)
		}

		err = toml.Unmarshal(mysteryGiftData, &mysteryGift)

		if err != nil {
			panic(err)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	// Checks if the player can recieve the gift based off the system clock(can be modified at any time by the user)
	beginningDate, _ := time.Parse("2006-01-02", mysteryGift.BeginningDate)
	endDate, _ := time.Parse("2006-01-02", mysteryGift.EndDate)
	playerDate, _ := time.Parse("2006-01-02", currentDate)
	if TimeCheck(beginningDate, endDate, playerDate) {
		return mysteryGift
	}

	return mysteryGift
}

func main() {
	SetUpDB()
	fmt.Print("Set up db!\n")
	InsertDBData()
}
