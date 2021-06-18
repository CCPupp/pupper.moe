package discord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"secret"

	"github.com/CCPupp/states.osutools/internal/commands"
	"github.com/CCPupp/states.osutools/internal/player"
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

	os.Exit(0)

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || len(m.Content) < 1 {
		return
	}
	// If the message starts with the prefix, handle accordingly.
	if m.Content[0:1] == prefix {
		command := m.Content[1:]
		fmt.Println(m.Author.Username + " used command: `" + command + "`")
		parts := strings.Split(command, " ")
		length := len(parts)
		if parts[0] == "ping" {
			s.ChannelMessageSend(m.ChannelID, commands.Ping())
		} else if parts[0] == "getuser" || parts[0] == "stats" {
			if len(parts) == 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, commands.GetUser(strconv.Itoa(player.GetUserByDiscordId(m.Author.ID).ID)))
			} else {
				if strings.HasPrefix(parts[1], "<@!") {
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetUser(strconv.Itoa(player.GetUserByDiscordId(strings.Trim(parts[1], "<@!>")).ID)))
				} else {
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetUser(parts[1]))
				}
			}
			// There 100% is a better way of writing this but I think this covers all possible inputs
		} else if parts[0] == "state" || parts[0] == "leaderboard" || parts[0] == "lb" {
			if len(parts) == 1 {
				// print page 1 of user's state
				s.ChannelMessageSendEmbed(m.ChannelID, commands.GetStateLeaderboard(player.GetUserByDiscordId(m.Author.ID).State, 1))
			} else if len(parts) == 2 {
				if page, err := strconv.Atoi(parts[1]); err == nil {
					// print page X of user's state
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetStateLeaderboard(player.GetUserByDiscordId(m.Author.ID).State, page))
				} else {
					// print page 1 of given state
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetStateLeaderboard(parts[1], 1))
				}
			} else if len(parts) == 3 {
				if page, err := strconv.Atoi(parts[2]); err == nil {
					// print page X of given state
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetStateLeaderboard(parts[1], page))
				} else {
					// print page 1 of given state with a space
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetStateLeaderboard(parts[1]+" "+parts[2], 1))
				}
			} else if len(parts) == 4 {
				if page, err := strconv.Atoi(parts[3]); err == nil {
					// print page X of given state
					s.ChannelMessageSendEmbed(m.ChannelID, commands.GetStateLeaderboard(parts[1]+" "+parts[2], page))
				}
			}
		} else if parts[0] == "setadmin" && length > 1 {
			if m.Author.ID == adminID {
				idInt, _ := strconv.Atoi(parts[1])
				s.ChannelMessageSend(m.ChannelID, commands.AssignAdmin(player.GetUserById(idInt)))
			} else {
				s.ChannelMessageSend(m.ChannelID, "Only ponpar can run this command.")
			}
		} else if parts[0] == "link" && length > 1 {
			idInt, err := strconv.Atoi(parts[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Not a Valid osu! ID.")
			} else {
				s.ChannelMessageSend(m.ChannelID, commands.LinkDiscordAccount(player.GetUserById(idInt), m.Author))
			}
		} else if parts[0] == "help" {
			s.ChannelMessageSendEmbed(m.ChannelID, commands.Help())
		} else if parts[0] == "dump" {
			author := player.GetUserByDiscordId(m.Author.ID)
			if author.Admin {
				s.ChannelMessageSendComplex(m.ChannelID, commands.Dump())
			} else {
				s.ChannelMessageSend(m.ChannelID, "Only admins can run this command.")
			}
		}

	}

}
