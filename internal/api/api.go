package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"secret"
	"time"

	"github.com/CCPupp/pupper.moe/internal/player"
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

func GetUser(id string, token string) player.User {
	url := "https://osu.ppy.sh/api/v2/users/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	osuClient := http.Client{
		Timeout: time.Second * 5, // Timeout after 5 seconds
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
	log.Print("User: " + user.Username)
	return user
}

func GetMe(gamemode string, w http.ResponseWriter, r *http.Request, token string) player.User {
	url := "https://osu.ppy.sh/api/v2/me/osu"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	osuClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
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
