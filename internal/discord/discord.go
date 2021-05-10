package discord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"secret"

	"github.com/CCPupp/pupper.moe/internal/commands"
	"github.com/CCPupp/pupper.moe/internal/player"
	"github.com/bwmarrin/discordgo"
)

type Discords struct {
	Discords []Discord `json:"discords"`
}

type Discord struct {
	State string `json:"state"`
	Link  string `json:"link"`
}

const prefix = "-"
const botToken = secret.DISCORD_TOKEN

//ponpar discord ID
const adminID = "98190856254676992"

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

func StartBot() {
	dg, err := discordgo.New("Bot " + botToken)

	if err != nil {
		fmt.Println(err)
		return
	}

	dg.AddHandler(messageCreate)
	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message starts with the prefix, handle accordingly.
	if m.Content[0:1] == prefix {
		command := m.Content[1:]
		if command == "ping" {
			s.ChannelMessageSend(m.ChannelID, commands.Ping(command))
		} else if command[0:7] == "getuser" {
			s.ChannelMessageSendEmbed(m.ChannelID, commands.GetUser(command[8:]))
		} else if command[0:8] == "setadmin" {
			if m.Author.ID == adminID {
				idInt, _ := strconv.Atoi(command[9:])
				s.ChannelMessageSend(m.ChannelID, commands.AssignAdmin(player.GetUserById(idInt)))
			} else {
				s.ChannelMessageSend(m.ChannelID, "Only ponpar can run this command.")
			}
		} else if command[0:4] == "link" {
			idInt, _ := strconv.Atoi(command[5:])
			s.ChannelMessageSend(m.ChannelID, commands.LinkDiscordAccount(player.GetUserById(idInt), m.Author))
		}

	}

}
