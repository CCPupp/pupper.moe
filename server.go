// websockets.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"secret"

	"github.com/CCPupp/pupper.moe/internal/api"
	"github.com/CCPupp/pupper.moe/internal/htmlbuilder"
	"github.com/CCPupp/pupper.moe/internal/player"
	"github.com/CCPupp/pupper.moe/internal/updater"
	"github.com/CCPupp/pupper.moe/internal/validations"

	_ "github.com/bmizerany/pq"
)

func main() {
	go updater.StartUpdate()
	// Handler points to available directories
	http.Handle("/web/html", http.StripPrefix("/web/html", http.FileServer(http.Dir("web/html"))))
	http.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", http.FileServer(http.Dir("web/scripts"))))
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css"))))
	http.Handle("/web/media/", http.StripPrefix("/web/media/", http.FileServer(http.Dir("web/media"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "web/html/index.html")
		} else if r.URL.Path[1:4] == "all" {
			fmt.Fprint(w, htmlbuilder.CreateAllHTML())
		} else if r.URL.Path[1:5] == "user" {
			fmt.Fprint(w, user(w, r))
		} else if r.URL.Path[1:6] == "login" {
			http.Redirect(w, r, "https://osu.ppy.sh/oauth/authorize?response_type=code&client_id="+strconv.Itoa(secret.OSU_CLIENT_ID)+"&redirect_uri="+secret.REDIRECT_URL+"/user&scope=public", http.StatusSeeOther)
		} else if r.URL.Path[1:7] == "states" {
			if r.URL.Path[8:] != "" {
				fmt.Fprint(w, htmlbuilder.CreateStateHTML(r.URL.Path[8:]))
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
		//mode := r.FormValue("mode")
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
			//player.SetUserMode(mode, id)
			idInt, _ := strconv.Atoi(id)
			user := player.GetUserById(idInt)
			fmt.Fprint(w, (htmlbuilder.CreateUser(user)))
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
		if id != "0" {
			if validations.ValidateState(state) {
				if player.CheckStateLock(idInt) {
					fmt.Fprint(w, "<h2>This user is locked!</h2>")
				} else {
					createUserFromId(id, state)
					fmt.Fprint(w, "<h2>Success!</h2>")
				}
			} else {
				fmt.Fprint(w, "<h2>Invalid State.</h2>")
			}
		} else {
			fmt.Fprint(w, "<h2>User Already Exists.</h2>")
		}

	})

	//Serves local webpage for testing
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

func user(w http.ResponseWriter, r *http.Request) string {
	token := api.GetUserToken(r.URL.Query().Get("code"))
	id := api.GetMe("osu", w, r, token)
	user := api.GetUser(strconv.Itoa(id.ID), token)
	if player.CheckDuplicate(user.ID) {
		player.OverwriteExisting(player.RetrieveUser(user.ID), user)
	} else {
		player.WriteToUser(user)
	}
	finalString := htmlbuilder.BuildHTMLHeader()
	finalString += htmlbuilder.BuildHTMLNavbar()
	finalString += htmlbuilder.CreateUser(player.RetrieveUser(user.ID))
	finalString += htmlbuilder.CreateOptions(player.RetrieveUser(user.ID))
	finalString += htmlbuilder.BuildHTMLFooter()
	return finalString
}

func createUserFromId(id string, state string) {
	clientToken := api.GetClientToken()
	var newUser player.User = api.GetUser(id, clientToken)
	newUser.State = state
	player.WriteToUser(newUser)
}
