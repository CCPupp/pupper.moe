package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"secret"
	"strconv"
	"time"

	"states.osutools/achievement"
	"states.osutools/player"
)

type UserRequest struct {
	Grant_type    string `json:"grant_type"`
	Client_id     int    `json:"client_id"`
	Client_secret string `json:"client_secret"`
	Redirect_uri  string `json:"redirect_uri"`
	Code          string `json:"code"`
}

type ClientRequest struct {
	Grant_type    string `json:"grant_type"`
	Client_id     int    `json:"client_id"`
	Client_secret string `json:"client_secret"`
	Scope         string `json:"scope"`
}

type Token struct {
	Token_type    string `json:"token_type"`
	Expires_in    int    `json:"expires_in"`
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
}

// GetRecent returns 10 recent plays from given user ID
func GetRecent(id int, token string) []achievement.Score {
	url := "https://osu.ppy.sh/api/v2/users/" + (strconv.Itoa(id)) + "/scores/recent?include_fails=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	osuClient := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, getErr := osuClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var event []achievement.Score
	jsonErr := json.Unmarshal(body, &event)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return event
}

// GetUser returns User with data from the osu! APIv2
func GetUser(id, token string) player.User {
	url := "https://osu.ppy.sh/api/v2/users/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	osuClient := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, getErr := osuClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var user player.User
	jsonErr := json.Unmarshal(body, &user)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return user
}

// GetMe returns User with data from the osu! APIv2 using their default game mode
func GetMe(token string) player.User {
	url := "https://osu.ppy.sh/api/v2/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	osuClient := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, getErr := osuClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var user player.User
	jsonErr := json.Unmarshal(body, &user)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return user
}

func GetUserToken(key string) string {
	url := "https://osu.ppy.sh/oauth/token"
	var jsonStr, _ = json.Marshal(UserRequest{
		Grant_type:    "authorization_code",
		Client_id:     secret.OSU_CLIENT_ID,
		Client_secret: secret.OSU_CLIENT_SECRET,
		Redirect_uri:  secret.REDIRECT_URL + "/user",
		Code:          key})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var token Token
	jsonErr := json.Unmarshal(body, &token)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return token.Access_token
}

func GetClientToken() string {
	url := "https://osu.ppy.sh/oauth/token"
	var jsonStr, _ = json.Marshal(ClientRequest{
		Client_id:     secret.OSU_CLIENT_ID,
		Client_secret: secret.OSU_CLIENT_SECRET,
		Grant_type:    "client_credentials",
		Scope:         "public"})
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var token Token
	jsonErr := json.Unmarshal(body, &token)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return token.Access_token
}
