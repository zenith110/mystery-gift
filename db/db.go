package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dgraph-io/badger/v4"
	Models "github.com/zenith110/pokemon-go-engine-toml-models/models"
)

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
		panic(err)
	}
	for mysteryGiftIndex := range mysteryGifts.Mysterygifts {
		err = client.Update(func(txn *badger.Txn) error {
			if mysteryGifts.Mysterygifts[mysteryGiftIndex].GiftType == "Pokemon" {
				data := fmt.Sprintf("%v", mysteryGifts.Mysterygifts[mysteryGiftIndex].Pokemon)
				txn.Set([]byte(mysteryGifts.Mysterygifts[mysteryGiftIndex].Name), []byte(data))
				return nil
			}
			return nil
		})
		if err != nil {
			panic(err)
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
