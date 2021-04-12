package discord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Discords struct {
	Discords []Discord `json:"discords"`
}

type Discord struct {
	State string `json:"state"`
	Link  string `json:"link"`
}

func GetDiscordJSON() Discords {
	// Open our jsonFile
	discordJsonFile, err := os.Open("web/data/discords.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer discordJsonFile.Close()

	discordByteValue, _ := ioutil.ReadAll(discordJsonFile)

	var discords Discords

	json.Unmarshal(discordByteValue, &discords)

	return discords
}
