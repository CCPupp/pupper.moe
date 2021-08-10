package stats

import (
	"time"

	"github.com/CCPupp/states.osutools/internal/player"
)

var StatsRunning = false
var TotalUsers = 0

func StartStats() {
	for {
		StatsRunning = true
		setTotalUsers()
		time.Sleep(time.Second * 1)
	}
}

func setTotalUsers() {
	TotalUsers = len(player.UserList)
}
