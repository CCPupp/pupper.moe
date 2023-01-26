package updater

import (
	"fmt"
	"strconv"
	"time"

	"states.osutools/api"
	"states.osutools/player"
)

var IsUpdating = false
var Progress = 0.0

func StartUpdate() {
	go updateLoop()
}

func updateLoop() {
	for {
		go doUpdate()
		time.Sleep(1 * time.Hour)
	}
}

func doUpdate() {
	IsUpdating = true
	fmt.Println("Starting Update")
	go progressLoop()
	clientToken := api.GetClientToken()
	for i := 0; i < len(player.UserList); i++ {
		Progress = (float64(i) / float64(len(player.UserList)))
		updateUser(player.UserList, i, clientToken)
	}
	IsUpdating = false
	fmt.Println("Update Complete!")
	player.BackupUserJSON()
	fmt.Println("Backups Created!")
}

func updateUser(users []player.User, i int, clientToken string) {
	updatedUser, err := api.GetUser(strconv.Itoa(users[i].ID), clientToken)
	if err != nil {
		player.OverwriteExistingUser(player.GetUserById(users[i].ID), updatedUser)
	}
}

func progressLoop() {
	for IsUpdating {
		reportProgress()
		time.Sleep(15 * time.Second)
	}
}

func reportProgress() {
	fmt.Print(Progress * 100)
	fmt.Println("%")
}
