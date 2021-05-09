package commands

import (
	"strconv"

	"github.com/CCPupp/pupper.moe/internal/player"
	"github.com/bwmarrin/discordgo"
)

func Ping(content string) string {
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

func makeUserFields(user player.User) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}
	mode := discordgo.MessageEmbedField{
		Name:   "Mode",
		Value:  user.Playmode,
		Inline: false,
	}
	state := discordgo.MessageEmbedField{
		Name:   "State",
		Value:  user.State,
		Inline: false,
	}
	stateRank := discordgo.MessageEmbedField{
		Name:   "State Rank",
		Value:  strconv.Itoa(player.GetUserStateRank(user.ID, user.State)),
		Inline: false,
	}
	globalRank := discordgo.MessageEmbedField{
		Name:   "Global Rank",
		Value:  strconv.Itoa(user.Statistics.Global_rank),
		Inline: false,
	}
	fields = append(fields, &mode, &state, &stateRank, &globalRank)
	// finalString := `
	// 	**Mode:** ` + user.Playmode +
	// 	`\n**State:** ` + user.State +
	// 	`\n**State Rank:** ` + strconv.Itoa(player.GetUserStateRank(user.ID, user.State)) +
	// 	`\n**Global Rank:** ` + strconv.Itoa(user.Statistics.Global_rank)
	return fields
}
