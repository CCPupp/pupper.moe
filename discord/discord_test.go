package discord

import "testing"

func TestInitializeDiscords(t *testing.T) {
	var DiscordList []Discord

	InitializeDiscords()

	if DiscordList == nil {
		t.Error("Initialization FAILED")
	} else {
		t.Log("Initialization PASSED")
	}
}
