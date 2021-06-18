package updater

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CCPupp/states.osutools/internal/api"
	"github.com/CCPupp/states.osutools/internal/player"
)

var IsUpdating = false

func StartUpdate() {
	for {
		doUpdate()
		time.Sleep(8 * time.Hour)
	}
}

func doUpdate() {

	IsUpdating = true
	fmt.Println("Starting Update")
	clientToken := api.GetClientToken()
	users := player.GetUserJSON()
	for i := 0; i < len(users.Users); i++ {
		time.Sleep(300 * time.Millisecond)
		updateUser(users, i, clientToken)
	}
	player.BackupUserJSON()
	IsUpdating = false
	fmt.Println("Update Complete!")

}

func updateUser(users player.Users, i int, clientToken string) {
	updatedUser := api.GetUser(strconv.Itoa(users.Users[i].ID), clientToken)
	player.OverwriteExistingUser(player.GetUserById(users.Users[i].ID), updatedUser)
}
