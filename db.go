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

func (mgdb MysteryGiftDB) SetUpDB() *badger.DB {
	if client != nil {
		var err error
		opt := badger.DefaultOptions("").WithInMemory(true)
		client, err = badger.Open(opt)
		if err != nil {
			fmt.Printf("Error while setting up db!\nError is %v", err)
		}
		return client
	}
	return client
}

func (mgdb MysteryGiftDB) InsertDBData() {
	file, err := os.Open("gifts.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var mysteryGifts []Models.MysteryGift

	err = toml.Unmarshal(bytes, &mysteryGifts)

	if err != nil {
		panic(err)
	}
	for mysteryGiftIndex := range mysteryGifts {
		err = client.Update(func(txn *badger.Txn) error {
			if mysteryGifts[mysteryGiftIndex].GiftType == "Pokemon" {
				data := fmt.Sprintf("%v", mysteryGifts[mysteryGiftIndex].Pokemon)
				txn.Set([]byte(mysteryGifts[mysteryGiftIndex].Name), []byte(data))
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
func (mgdb MysteryGiftDB) SearchDBData(currentDate string, giftName string) Models.MysteryGift {
	var mysteryGiftData []byte
	var mysteryGift Models.MysteryGift
	err := client.View(func(txn *badger.Txn) error {
		data, _ := txn.Get([]byte(giftName))
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
	beginningDate, _ := time.Parse("2006-01-02", mysteryGift.BeginningDate)
	endDate, _ := time.Parse("2006-01-02", mysteryGift.EndDate)
	playerDate, _ := time.Parse("2006-01-02", currentDate)
	if TimeCheck(beginningDate, endDate, playerDate) {
		return mysteryGift
	}
	return mysteryGift
}
