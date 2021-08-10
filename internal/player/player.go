package player

// Player stores information about the player to parse onto the webpage
type User struct {
	Statistics     Statistic `json:"statistics"`
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	ProfileColor   string    `json:"profile_colour"`
	AvatarURL      string    `json:"avatar_url"`
	Discord        string    `json:"discord"`
	CoverURL       string    `json:"cover_url"`
	CountryCode    string    `json:"country_code"`
	Playmode       string    `json:"playmode"`
	ReplaysWatched int       `json:"replays_watched_by_others"`
	Badges         []Badge   `json:"badges"`
	// These items are not pulled from the osu!api and instead are stored locally.
	State      string    `json:"state"`
	AdvState   string    `json:"advstate"`
	Background string    `json:"background"`
	Locks      Lock_info `json:"locks"`
	Admin      bool      `json:"admin"`
	DiscordID  string    `json:"discordid"`
}

type Statistic struct {
	Pp          float64    `json:"pp"`
	Global_rank int        `json:"Global_rank"`
	Accuracy    float64    `json:"hit_accuracy"`
	Play_count  int        `json:"play_count"`
	Level       Level_info `json:"level"`
}

type Level_info struct {
	Current  int `json:"current"`
	Progress int `json:"progress"`
}

type Lock_info struct {
	Mode_Lock  bool `json:"modelock"`
	State_Lock bool `json:"statelock"`
}

type Users struct {
	Users []User `json:"users"`
}

type Badge struct {
	Awarded_At  string `json:"awarded_at" `
	Description string `json:"description" `
	Image_URL   string `json:"image_url"  `
	URL         string `json:"url"`
}

type Event struct {
	CreatedAt string `json:"created_at"`
	Id        int    `json:"id"`
	Type      string `json:"type"`
}

type RankEvent struct {
	CreatedAt string `json:"created_at"`
	Id        int    `json:"id"`
	Type      string `json:"type"`
	// TODO: FIX WHEN PEPPY FIXES
	ScoreRank string       `json:"scoreRank"`
	Rank      int          `json:"rank"`
	Mode      string       `json:"mode"`
	Beatmap   BeatmapEvent `json:"Beatmap"`
	User      UserEvent    `json:"User"`
}

type UserEvent struct {
	CreatedAt string `json:"created_at"`
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	URL       string `json:"url"`
}

type BeatmapEvent struct {
	CreatedAt string `json:"created_at"`
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	URL       string `json:"url"`
}
