package updater

import (
	"fmt"
	"strconv"
	"time"

	"states.osutools/api"
	"states.osutools/player"
)

var IsUpdating = false

func StartUpdate() {
	go updateLoop()
	go backupLoop()
}

func updateLoop() {
	for {
		go doUpdate()
		time.Sleep(8 * time.Hour)
	}
}

func doUpdate() {
	IsUpdating = true
	fmt.Println("Starting Update")
	clientToken := api.GetClientToken()
	for i := 0; i < len(player.UserList); i++ {
		updateUser(player.UserList, i, clientToken)
	}
	IsUpdating = false
	fmt.Println("Update Complete!")
}

func updateUser(users []player.User, i int, clientToken string) {
	updatedUser := api.GetUser(strconv.Itoa(users[i].ID), clientToken)
	player.NewOverwriteExistingUser(player.NewGetUserById(users[i].ID), updatedUser)
}

func backupLoop() {
	for {
		go backup()
		time.Sleep(30 * time.Minute)
	}
}

func backup() {
	player.NewBackupUserJSON()
}
