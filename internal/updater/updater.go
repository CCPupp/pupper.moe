package updater

import (
	"fmt"
	"strconv"
	"time"

	"github.com/CCPupp/pupper.moe/internal/api"
	"github.com/CCPupp/pupper.moe/internal/player"
)

var IsUpdating = false

func StartUpdate() {
	for {
		doUpdate()
		time.Sleep(24 * time.Hour)
	}
}

func doUpdate() {

	IsUpdating = true
	fmt.Println("Starting Update")
	clientToken := api.GetClientToken()
	users := player.GetUserJSON()
	for i := 0; i < len(users.Users); i++ {
		time.Sleep(100 * time.Millisecond)
		updateUser(users, i, clientToken)
	}
	IsUpdating = false
	fmt.Println("Update Complete!")

}

func updateUser(users player.Users, i int, clientToken string) {
	updatedUser := api.GetUser(strconv.Itoa(users.Users[i].ID), clientToken)
	player.OverwriteExisting(player.RetrieveUser(users.Users[i].ID), updatedUser)
}
