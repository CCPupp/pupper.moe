// websockets.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"secret"

	"github.com/CCPupp/pupper.moe/internal/api"
	"github.com/CCPupp/pupper.moe/internal/discord"
	"github.com/CCPupp/pupper.moe/internal/htmlbuilder"
	"github.com/CCPupp/pupper.moe/internal/player"
	"github.com/CCPupp/pupper.moe/internal/updater"
	"github.com/CCPupp/pupper.moe/internal/validations"

	_ "github.com/bmizerany/pq"
)

const OnlyBot = false

func main() {
	if !secret.IS_TESTING {
		go updater.StartUpdate()
	}
	go discord.StartBot()
	// Handler points to available directories
	http.Handle("/web/html", http.StripPrefix("/web/html", http.FileServer(http.Dir("web/html"))))
	http.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", http.FileServer(http.Dir("web/scripts"))))
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css"))))
	http.Handle("/web/media/", http.StripPrefix("/web/media/", http.FileServer(http.Dir("web/media"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "web/html/index.html")
		} else if r.URL.Path[1:4] == "all" {
			fmt.Fprint(w, htmlbuilder.CreateAllHTML(1))
		} else if r.URL.Path[1:5] == "user" {
			fmt.Fprint(w, user(r))
		} else if r.URL.Path[1:6] == "login" {
			http.Redirect(w, r, "https://osu.ppy.sh/oauth/authorize?response_type=code&client_id="+strconv.Itoa(secret.OSU_CLIENT_ID)+"&redirect_uri="+secret.REDIRECT_URL+"/user&scope=public", http.StatusSeeOther)
		} else if r.URL.Path[1:7] == "states" {
			if r.URL.Path[8:] != "" {
				if r.URL.Path[len(r.URL.Path)-3:] == "osu" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-4], "osu", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "mania" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-6], "mania", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "catch" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-6], "fruits", 2)
				} else if r.URL.Path[len(r.URL.Path)-5:] == "taiko" {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:len(r.URL.Path)-6], "taiko", 2)
				} else {
					htmlbuilder.CreateStateHTML(w, r.URL.Path[8:], "all", 1)
				}
			}
		} else {
			http.ServeFile(w, r, "web/html/"+r.URL.Path[1:]+".html")
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
		}
		state := r.FormValue("state")
		bg := r.FormValue("bg")
		id := r.FormValue("id")
		if id != "0" {
			if bg == "true" || bg == "false" {
				player.SetUserBg(bg, id)
			}
			if validations.ValidateState(state) {
				player.SetUserState(state, id)
			} else {
				fmt.Fprint(w, "<h2>Invalid State.</h2>")
			}
			idInt, _ := strconv.Atoi(id)
			user := player.GetUserById(idInt)
			fmt.Fprint(w, (htmlbuilder.CreateUser(user, 0)))
		} else {
			fmt.Fprint(w, "<h2>Submission Failed.</h2>")
		}

	})

	http.HandleFunc("/submitPlayer", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
		}
		state := r.FormValue("state")
		id := r.FormValue("id")
		idInt, _ := strconv.Atoi(id)
		if player.GetUserById(idInt).ID != idInt {
			if validations.ValidateState(state) {
				if player.CheckStateLock(idInt) {
					fmt.Fprint(w, "<h2>This user is locked!</h2>")
				} else {
					if updater.IsUpdating {
						fmt.Fprint(w, "<h2>Update in progress, please wait.</h2>")
					} else {
						createUserFromId(id, state)
						fmt.Fprint(w, "<h2>Success!</h2>")
					}
				}
			} else {
				fmt.Fprint(w, "<h2>Invalid State.</h2>")
			}
		} else {
			fmt.Fprint(w, "<h2>User Already Exists.</h2>")
		}

	})

	//Serves local webpage for testing
	if OnlyBot == false {
		if secret.IS_TESTING {
			errhttp := http.ListenAndServe(":8080", nil)
			if errhttp != nil {
				log.Fatal("Web server (HTTP): ", errhttp)
			}
		} else {
			//Serves the webpage to the internet
			errhttps := http.ListenAndServeTLS(":443", "certs/cert.pem", "certs/key.pem", nil)
			if errhttps != nil {
				log.Fatal("Web server (HTTPS): ", errhttps)
			}
		}
	}

}

func user(r *http.Request) string {
	token := api.GetUserToken(r.URL.Query().Get("code"))
	user := api.GetMe(token)
	var localUser player.User
	if player.CheckDuplicate(user.ID) {
		localUser = player.GetUserById(user.ID)
		player.OverwriteExistingUser(localUser, user)
	} else {
		player.WriteToUser(user)
		localUser = player.GetUserById(user.ID)
	}
	finalString := htmlbuilder.BuildHTMLHeader(1, "Just "+localUser.Username)
	finalString += htmlbuilder.BuildHTMLNavbar()
	finalString += htmlbuilder.CreateUser(localUser, 0)
	finalString += htmlbuilder.CreateOptions(localUser)
	if localUser.Admin {
		finalString += htmlbuilder.CreateAdminPanel(localUser)
	}
	finalString += htmlbuilder.BuildHTMLFooter()
	return finalString
}

func createUserFromId(id string, state string) {
	clientToken := api.GetClientToken()
	var newUser player.User = api.GetUser(id, clientToken)
	newUser.State = state
	player.WriteToUser(newUser)
}
