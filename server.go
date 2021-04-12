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

	_ "github.com/bmizerany/pq"
)

func main() {

	// Handler points to available directories
	http.Handle("/web/html", http.StripPrefix("/web/html", http.FileServer(http.Dir("web/html"))))
	http.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", http.FileServer(http.Dir("web/scripts"))))
	http.Handle("/web/css/", http.StripPrefix("/web/css/", http.FileServer(http.Dir("web/css"))))
	http.Handle("/web/media/", http.StripPrefix("/web/media/", http.FileServer(http.Dir("web/media"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[1:] == "" {
			http.ServeFile(w, r, "web/html/index.html")
		} else if r.URL.Path[1:4] == "all" {
			fmt.Fprintf(w, htmlbuilder.CreateAllHTML())
		} else if r.URL.Path[1:5] == "user" {
			fmt.Fprintf(w, user(w, r))
		} else if r.URL.Path[1:6] == "login" {
			http.Redirect(w, r, "https://osu.ppy.sh/oauth/authorize?response_type=code&client_id="+strconv.Itoa(secret.OSU_CLIENT_ID)+"&redirect_uri="+secret.REDIRECT_URL+"/user&scope=public", http.StatusSeeOther)
		} else if r.URL.Path[1:7] == "states" {
			fmt.Fprintf(w, htmlbuilder.CreateStateHTML(r.URL.Path[8:]))
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
		mode := r.FormValue("mode")
		id := r.FormValue("id")
		if bg == "true" || bg == "false" {
			player.SetUserBg(bg, id)
		}
		player.SetUserState(state, id)
		player.SetUserMode(mode, id)

		fmt.Fprintf(w, "<h2>Update Successful!</h2>")
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
	token := api.GetToken(r.URL.Query().Get("code"))
	id := api.GetMe("osu", w, r, token)
	user := api.GetUser(strconv.Itoa(id.ID), "osu", w, r, token)
	if player.CheckDuplicate(user.ID) {
		player.OverwriteExisting(player.RetrieveUser(user.ID), user)
	} else {
		player.WriteToUser(player.RetrieveUser(user.ID))
	}
	finalString := htmlbuilder.BuildHTMLHeader()
	finalString += htmlbuilder.BuildHTMLNavbar()
	finalString += htmlbuilder.CreateUser(player.RetrieveUser(user.ID))
	finalString += htmlbuilder.CreateOptions(player.RetrieveUser(user.ID))
	finalString += htmlbuilder.BuildHTMLFooter()
	return finalString
}
