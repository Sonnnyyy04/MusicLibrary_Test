package models

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"groupName"`
	SongName    string `json:"songName"`
	DateRelease string `json:"dateRelease"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
