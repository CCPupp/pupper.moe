package commands

import (
	"strconv"

	"github.com/CCPupp/pupper.moe/internal/player"
	"github.com/CCPupp/pupper.moe/internal/validations"
	"github.com/bwmarrin/discordgo"
)

func Ping() string {
	return "Pong!"
}

func Help() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "List of Commands:",
		Description: "[] notates optional fields, () are required fields, | notates an alternative usage for the same command",
		Fields:      makeHelpFields(),
	}
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
	if user.ID != 0 && user.Discord == discordUser.Username+"#"+discordUser.Discriminator {
		player.SetUserDiscordID(user, discordUser.ID)
		return user.Username + " is linked to " + discordUser.Mention() + "."
	} else {
		return "Invalid ID / ID not on userpage."
	}
}

func GetStateLeaderboard(state string) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title: "Invalid State / Account Not Linked",
	}
	if validations.ValidateState(state) {
		embed = discordgo.MessageEmbed{
			Title:  state,
			Fields: makeStateFields(state),
		}
		return &embed
	}
	return &embed
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
	badges := discordgo.MessageEmbedField{
		Name:   "BWS Rank",
		Value:  strconv.Itoa(player.GetBWSRank(user)),
		Inline: true,
	}
	fields = append(fields, &mode, &state, &stateRank, &globalRank, &badges)
	return fields
}

func makeStateFields(state string) []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}
	users := player.GetUserJSON()
	users = player.SortUsers(users)
	player1 := discordgo.MessageEmbedField{}
	player2 := discordgo.MessageEmbedField{}
	player3 := discordgo.MessageEmbedField{}
	player4 := discordgo.MessageEmbedField{}
	player5 := discordgo.MessageEmbedField{}
	count := 0
	errorEmbed := discordgo.MessageEmbedField{
		Name:   "Error",
		Value:  "Something went wrong.",
		Inline: true,
	}
	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].State == state {
			count++
			if count == 1 {
				player1 = discordgo.MessageEmbedField{
					Name:   strconv.Itoa(count),
					Value:  users.Users[i].Username,
					Inline: false,
				}
			}
			if count == 2 {
				player2 = discordgo.MessageEmbedField{
					Name:   strconv.Itoa(count),
					Value:  users.Users[i].Username,
					Inline: false,
				}
			}
			if count == 3 {
				player3 = discordgo.MessageEmbedField{
					Name:   strconv.Itoa(count),
					Value:  users.Users[i].Username,
					Inline: false,
				}
			}
			if count == 4 {
				player4 = discordgo.MessageEmbedField{
					Name:   strconv.Itoa(count),
					Value:  users.Users[i].Username,
					Inline: false,
				}
			}
			if count == 5 {
				player5 = discordgo.MessageEmbedField{
					Name:   strconv.Itoa(count),
					Value:  users.Users[i].Username,
					Inline: false,
				}
			}
		}
	}

	if count == 0 {
		fields = append(fields, &errorEmbed)
	} else {
		fields = append(fields, &player1, &player2, &player3, &player4, &player5)
	}
	return fields
}

func makeHelpFields() []*discordgo.MessageEmbedField {
	fields := []*discordgo.MessageEmbedField{}
	ping := discordgo.MessageEmbedField{
		Name:   "-ping",
		Value:  "Pong! A good way to check if the website & bot are online.",
		Inline: false,
	}
	stats := discordgo.MessageEmbedField{
		Name:   "-stats [@user | Discord ID | osu! ID] | -getuser",
		Value:  "The stats command shows information about the user from my database. If the user has connected their discord account using -link they can be pulled by @ing them.",
		Inline: false,
	}
	link := discordgo.MessageEmbedField{
		Name:   "-link (osu ID)",
		Value:  "Connects your discord account to your account on the website and verifies it with a checkmark. YOU MUST HAVE YOUR ACCOUNT ON YOUR OSU USERPAGE TO LINK.",
		Inline: false,
	}
	state := discordgo.MessageEmbedField{
		Name:   "-state [State] | -leaderboard | -lb",
		Value:  "Shows the top 5 users from the specified state, if no state is specified your linked accounts state is used.",
		Inline: false,
	}
	fields = append(fields, &ping, &stats, &link, &state)
	return fields
}
