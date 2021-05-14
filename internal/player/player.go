package player

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Player stores information about the player to parse onto the webpage
type User struct {
	Statistics     Statistic `json:"statistics"`
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	ProfileColor   string    `json:"profile_colour"`
	AvatarURL      string    `json:"avatar_url"`
	Discord        string    `json:"discord"`
	CoverURL       string    `json:"cover_url"`
	CountryCode    string    `json:"country_code"`
	Playmode       string    `json:"playmode"`
	ReplaysWatched int       `json:"replays_watched_by_others"`
	Badges         []Badge   `json:"badges"`
	// These items are not pulled from the osu!api and instead are stored locally.
	State      string    `json:"state"`
	Background string    `json:"background"`
	Locks      Lock_info `json:"locks"`
	Admin      bool      `json:"admin"`
	DiscordID  string    `json:"discordid"`
}

type Statistic struct {
	Pp          float64    `json:"pp"`
	Global_rank int        `json:"Global_rank"`
	Accuracy    float64    `json:"hit_accuracy"`
	Play_count  int        `json:"play_count"`
	Level       Level_info `json:"level"`
}

type Level_info struct {
	Current  int `json:"current"`
	Progress int `json:"progress"`
}

type Lock_info struct {
	Mode_Lock  bool `json:"modelock"`
	State_Lock bool `json:"statelock"`
}

type Users struct {
	Users []User `json:"users"`
}

type Badge struct {
	Awarded_At  string `json:"awarded_at" `
	Description string `json:"description" `
	Image_URL   string `json:"image_url"  `
	URL         string `json:"url"`
}

func GetUserById(id int) User {
	users := GetUserJSON()

	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].ID == id {
			return users.Users[i]
		}
	}

	var nullUser User
	return nullUser

}

func GetUserByDiscordId(id string) User {
	users := GetUserJSON()

	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].DiscordID == id {
			return users.Users[i]
		}
	}

	var nullUser User
	return nullUser
}

func GetUserJSON() Users {
	jsonFile, err := os.Open("web/data/users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users Users

	json.Unmarshal(byteValue, &users)

	return users
}

func SortUsers(list Users) Users {
	sort.SliceStable(list.Users, func(i, j int) bool {
		return list.Users[i].Statistics.Global_rank < list.Users[j].Statistics.Global_rank
	})

	return list
}

func CheckDuplicate(dupe int) bool {
	currentList := GetUserJSON()

	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == dupe {
			return true
		}
	}

	return false
}

func GetBWSRank(user User) int {
	if user.Badges != nil {
		var total = 0
		for i := 0; i < len(user.Badges); i++ {
			total++
		}
		//Math.round(Math.pow(rank, Math.pow(0.9937, Math.pow(badges, 2)))).toLocaleString();
		bwsRankFloat := math.Round(math.Pow((float64(user.Statistics.Global_rank)), math.Pow(0.9937, math.Pow(float64(total), 2))))
		tempString := fmt.Sprintf("%f", bwsRankFloat)
		parts := strings.Split(tempString, ".")
		finalInt, _ := strconv.Atoi(parts[0])
		return finalInt
	}
	return user.Statistics.Global_rank
}

func CheckStateLock(id int) bool {
	user := GetUserById(id)
	return user.Locks.State_Lock
}

func SetUserState(state string, id string) {
	currentList := GetUserJSON()

	for i := 0; i < len(currentList.Users); i++ {
		if strconv.Itoa(currentList.Users[i].ID) == id {
			currentList.Users[i].State = state
			currentList.Users[i].Locks.State_Lock = true
		}
	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func SetUserBg(bg string, id string) {
	currentList := GetUserJSON()

	for i := 0; i < len(currentList.Users); i++ {
		if strconv.Itoa(currentList.Users[i].ID) == id {
			currentList.Users[i].Background = bg
		}
	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func SetUserMode(mode string, id string) {
	currentList := GetUserJSON()

	for i := 0; i < len(currentList.Users); i++ {
		if strconv.Itoa(currentList.Users[i].ID) == id {
			currentList.Users[i].Playmode = mode
			currentList.Users[i].Locks.Mode_Lock = true
		}
	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func SetUserAdmin(user User) {
	currentList := GetUserJSON()

	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == user.ID {
			currentList.Users[i].Admin = true
		}
	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func SetUserDiscordID(user User, discordID string) {
	currentList := GetUserJSON()

	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == user.ID {
			currentList.Users[i].DiscordID = discordID
		}
	}

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func WriteToUser(newUser User) {
	currentList := GetUserJSON()
	var badge Badge
	var badges []Badge
	for i := 0; i < len(newUser.Badges); i++ {
		badge = Badge{
			Awarded_At:  newUser.Badges[i].Awarded_At,
			Description: newUser.Badges[i].Description,
			Image_URL:   newUser.Badges[i].Image_URL,
			URL:         newUser.Badges[i].URL,
		}

		badges = append(badges, badge)
	}

	level := Level_info{
		Current:  newUser.Statistics.Level.Current,
		Progress: newUser.Statistics.Level.Progress,
	}

	stats := Statistic{
		Pp:          newUser.Statistics.Pp,
		Global_rank: newUser.Statistics.Global_rank,
		Accuracy:    newUser.Statistics.Accuracy,
		Play_count:  newUser.Statistics.Play_count,
		Level:       level,
	}

	currentList.Users = append(currentList.Users, User{
		ID:             newUser.ID,
		Username:       newUser.Username,
		CountryCode:    newUser.CountryCode,
		CoverURL:       newUser.CoverURL,
		Playmode:       newUser.Playmode,
		ProfileColor:   newUser.ProfileColor,
		AvatarURL:      newUser.AvatarURL,
		Discord:        newUser.Discord,
		ReplaysWatched: newUser.ReplaysWatched,
		Statistics:     stats,
		State:          newUser.State,
		Background:     newUser.Background,
		Locks:          newUser.Locks,
		Admin:          newUser.Admin,
		DiscordID:      newUser.DiscordID,
		Badges:         badges,
	})

	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func OverwriteExistingUser(existingUser User, pulledUser User) {
	currentList := GetUserJSON()
	var badge Badge
	var badges []Badge
	for i := 0; i < len(pulledUser.Badges); i++ {
		badge = Badge{
			Awarded_At:  pulledUser.Badges[i].Awarded_At,
			Description: pulledUser.Badges[i].Description,
			Image_URL:   pulledUser.Badges[i].Image_URL,
			URL:         pulledUser.Badges[i].URL,
		}

		badges = append(badges, badge)
	}

	level := Level_info{
		Current:  pulledUser.Statistics.Level.Current,
		Progress: pulledUser.Statistics.Level.Progress,
	}

	stats := Statistic{
		Pp:          pulledUser.Statistics.Pp,
		Global_rank: pulledUser.Statistics.Global_rank,
		Accuracy:    pulledUser.Statistics.Accuracy,
		Play_count:  pulledUser.Statistics.Play_count,
		Level:       level,
	}

	user := User{
		ID:             existingUser.ID,
		Username:       pulledUser.Username,
		CountryCode:    pulledUser.CountryCode,
		CoverURL:       pulledUser.CoverURL,
		Playmode:       pulledUser.Playmode,
		ProfileColor:   pulledUser.ProfileColor,
		AvatarURL:      pulledUser.AvatarURL,
		Discord:        pulledUser.Discord,
		ReplaysWatched: pulledUser.ReplaysWatched,
		Statistics:     stats,
		State:          existingUser.State,
		Background:     existingUser.Background,
		Locks:          existingUser.Locks,
		Admin:          existingUser.Admin,
		DiscordID:      existingUser.DiscordID,
		Badges:         badges,
	}

	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == existingUser.ID {
			currentList.Users[i] = user
		}
	}

	// now Marshal it
	finalList, _ := json.Marshal(currentList)

	ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func GetUserStateRank(id int, state string) int {
	users := SortUsers(GetUserJSON())
	rank := 0
	for i := 0; i < len(users.Users); i++ {
		if (users.Users[i].State == state) && (users.Users[i].Statistics.Global_rank != 0) {
			rank++
			if users.Users[i].ID == id {
				return rank
			}
		}
	}

	return 0
}

func GetTotalVerified() string {
	users := SortUsers(GetUserJSON())
	total := 0
	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].Locks.State_Lock {
			total++
		}
	}

	return strconv.Itoa(total)
}
