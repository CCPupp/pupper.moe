package discord

type Discords struct {
	Discords []Discord `json:"discords"`
}

type Discord struct {
	State string `json:"state"`
	Link  string `json:"link"`
}
