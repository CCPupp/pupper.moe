package player

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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
	State          string    `json:"state"`
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

type Users struct {
	Users []User `json:"users"`
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
	jsonFile, err := os.Open("web/data/users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var currentList Users
	json.Unmarshal(byteValue, &currentList)

	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == dupe {
			return true
		}
	}

	return false
}

func AddUserState(state string, user User) {
	jsonFile, err := os.Open("web/data/users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var currentList Users
	json.Unmarshal(byteValue, &currentList)

	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == user.ID {
			level := Level_info{
				Current:  user.Statistics.Level.Current,
				Progress: user.Statistics.Level.Progress,
			}

			stats := Statistic{
				Pp:          user.Statistics.Pp,
				Global_rank: user.Statistics.Global_rank,
				Accuracy:    user.Statistics.Accuracy,
				Play_count:  user.Statistics.Play_count,
				Level:       level,
			}

			currentList.Users = append(currentList.Users, User{
				ID:             user.ID,
				Username:       user.Username,
				State:          user.State,
				CountryCode:    user.CountryCode,
				CoverURL:       user.CoverURL,
				Playmode:       user.Playmode,
				ProfileColor:   user.ProfileColor,
				AvatarURL:      user.AvatarURL,
				Discord:        user.Discord,
				ReplaysWatched: user.ReplaysWatched,
				Statistics:     stats,
			})
		}
	}
}

func WriteToUser(newUser User) {
	jsonFile, err := os.Open("web/data/users.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var currentList Users
	json.Unmarshal(byteValue, &currentList)

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
		State:          newUser.State,
		CountryCode:    newUser.CountryCode,
		CoverURL:       newUser.CoverURL,
		Playmode:       newUser.Playmode,
		ProfileColor:   newUser.ProfileColor,
		AvatarURL:      newUser.AvatarURL,
		Discord:        newUser.Discord,
		ReplaysWatched: newUser.ReplaysWatched,
		Statistics:     stats,
	})

	finalList, _ := json.Marshal(currentList)

	err = ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func OverwriteExisting(existingUser User) {
	// Open our jsonFile
	jsonFile, err := os.Open("web/data/users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var currentList Users
	json.Unmarshal(byteValue, &currentList)
	// TODO: FIX THIS
	for i := 0; i < len(currentList.Users); i++ {
		if currentList.Users[i].ID == existingUser.ID {
		}
	}

	// now Marshal it
	finalList, _ := json.Marshal(currentList)

	err = ioutil.WriteFile("web/data/users.json", finalList, 0644)
}

func RetrieveUser(id int) User {
	jsonFile, err := os.Open("web/data/users.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Players array
	var users Users
	var user User

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'players' which we defined above
	json.Unmarshal(byteValue, &users)
	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].ID == id {
			user = users.Users[i]
		}
	}

	return user
}
