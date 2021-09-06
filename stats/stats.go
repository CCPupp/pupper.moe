package stats

import (
	"time"

	"states.osutools/player"
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
