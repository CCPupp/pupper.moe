package player

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"
)

var UserList []User

func InitializeUserList() {
	jsonFile, err := os.Open("web/data/users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users []User

	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		fmt.Println(err)
	}

	UserList = users
	clean()
	alphabetizeUsers()
}

func alphabetizeUsers() {
	var sortedList []User = UserList
	sort.SliceStable(sortedList, func(i, j int) bool {
		return UserList[i].Username > UserList[j].Username
	})
	UserList = sortedList
}

func clean() {
	for i := 0; i < len(UserList); i++ {
		if UserList[i].Username == "" || UserList[i].State == "" {
			DeleteUserById(UserList[i].ID)
		}
	}
}

func BackupUserJSON() {
	byteValue, _ := json.Marshal(UserList)
	ioutil.WriteFile("web/data/users.json", byteValue, 0644)
	ioutil.WriteFile("web/data/backups/usersBACKUP"+time.Now().String()+".json", byteValue, 0644)
}

func SortUsersByRank() []User {
	var sortedList []User = UserList
	sort.SliceStable(sortedList, func(i, j int) bool {
		return UserList[i].Statistics.Global_rank < UserList[j].Statistics.Global_rank
	})
	return sortedList
}

func DeleteUserById(id int) {
	var nullUser User
	for i := 0; i < len(UserList); i++ {
		if UserList[i].ID == id {
			UserList[i] = UserList[len(UserList)-1]
			UserList[len(UserList)-1] = nullUser
			UserList = UserList[:len(UserList)-1]
		}
	}
}

func GetUserById(id int) User {

	for i := 0; i < len(UserList); i++ {
		if UserList[i].ID == id {
			return UserList[i]
		}
	}

	var nullUser User
	return nullUser
}

func GetUserByDiscordId(id string) User {
	for i := 0; i < len(UserList); i++ {
		if UserList[i].DiscordID == id {
			return UserList[i]
		}
	}

	var nullUser User
	return nullUser
}

func CheckDuplicate(dupe int) bool {

	for i := 0; i < len(UserList); i++ {
		if UserList[i].ID == dupe {
			return true
		}
	}

	return false
}

func CheckStateLock(id int) bool {
	user := GetUserById(id)
	return user.Locks.State_Lock
}

func SetUserState(state string, id string, verified bool) {
	for i := 0; i < len(UserList); i++ {
		if strconv.Itoa(UserList[i].ID) == id {
			UserList[i].State = state
			UserList[i].Locks.State_Lock = verified
		}
	}
}

func SetUserAdvState(advstate string, id string) {
	for i := 0; i < len(UserList); i++ {
		if strconv.Itoa(UserList[i].ID) == id {
			UserList[i].AdvState = advstate
		}
	}
}

func SetUserBg(bg string, id string) {
	for i := 0; i < len(UserList); i++ {
		if strconv.Itoa(UserList[i].ID) == id {
			UserList[i].Background = bg
		}
	}
}

func SetUserMode(mode string, id string) {
	for i := 0; i < len(UserList); i++ {
		if strconv.Itoa(UserList[i].ID) == id {
			UserList[i].Playmode = mode
			UserList[i].Locks.Mode_Lock = true
		}
	}
}

func SetUserAdmin(user User) {
	for i := 0; i < len(UserList); i++ {
		if UserList[i].ID == user.ID {
			UserList[i].Admin = true
		}
	}
}

func SetUserDiscordID(user User, discordID string) {
	for i := 0; i < len(UserList); i++ {
		if UserList[i].ID == user.ID {
			UserList[i].DiscordID = discordID
		}
	}
}

func WriteToUser(newUser User) {
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

	UserList = append(UserList, User{
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
		AdvState:       newUser.AdvState,
		Background:     newUser.Background,
		Locks:          newUser.Locks,
		Admin:          newUser.Admin,
		DiscordID:      newUser.DiscordID,
		Badges:         badges,
	})
}
func OverwriteExistingUser(existingUser User, pulledUser User) {
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
		AdvState:       existingUser.AdvState,
		Background:     existingUser.Background,
		Locks:          existingUser.Locks,
		Admin:          existingUser.Admin,
		DiscordID:      existingUser.DiscordID,
		Badges:         badges,
	}

	for i := 0; i < len(UserList); i++ {
		if UserList[i].ID == existingUser.ID {
			UserList[i] = user
		}
	}
}

func GetUserStateRank(id int, state string) int {
	rank := 0
	for i := 0; i < len(UserList); i++ {
		if (UserList[i].State == state) && (UserList[i].Statistics.Global_rank != 0) {
			rank++
			if UserList[i].ID == id {
				return rank
			}
		}
	}

	return 0
}

func GetTotalVerified() string {
	total := 0
	for i := 0; i < len(UserList); i++ {
		if UserList[i].Locks.State_Lock {
			total++
		}
	}

	return strconv.Itoa(total)
}
