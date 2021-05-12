package commands

import (
	"strconv"

	"github.com/CCPupp/pupper.moe/internal/player"
	"github.com/bwmarrin/discordgo"
)

func Ping() string {
	return "Pong!"
}

func GetUser(id string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title: "Invalid ID",
	}
	if idInt, err := strconv.Atoi(id); err == nil {
		user := player.GetUserById(idInt)
		if user.ID != 0 {
			embed = discordgo.MessageEmbed{
				Title:  user.Username,
				Fields: makeUserFields(user),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL:    user.AvatarURL,
					Width:  15,
					Height: 15,
				},
			}
			return &embed
		}
	}
	return &embed
}

func AssignAdmin(user player.User) string {
	if user.ID != 0 {
		player.SetUserAdmin(user)
		return user.Username + " is now an admin."
	} else {
		return "Invalid ID"
	}
}

func LinkDiscordAccount(user player.User, discordUser *discordgo.User) string {
	if user.ID != 0 {
		player.SetUserDiscordID(user, discordUser.ID)
		return user.Username + " is linked to " + discordUser.Mention() + "."
	} else {
		return "Invalid ID"
	}
}

func makeUserFields(user player.User) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}
	mode := discordgo.MessageEmbedField{
		Name:   "Mode",
		Value:  user.Playmode,
		Inline: true,
	}
	state := discordgo.MessageEmbedField{
		Name:   "State",
		Value:  user.State,
		Inline: true,
	}
	stateRank := discordgo.MessageEmbedField{
		Name:   "State Rank",
		Value:  strconv.Itoa(player.GetUserStateRank(user.ID, user.State)),
		Inline: true,
	}
	globalRank := discordgo.MessageEmbedField{
		Name:   "Global Rank",
		Value:  strconv.Itoa(user.Statistics.Global_rank),
		Inline: true,
	}
	fields = append(fields, &mode, &state, &stateRank, &globalRank)
	return fields
}
