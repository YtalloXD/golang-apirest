package models

import "time"

// Game represents a video game in our system
type Game struct {
	ID          string    `json:"id"`
	GameName    string    `json:"game_name"`
	Publisher   string    `json:"publisher"`
	Developer   string    `json:"developer"`
	ReleaseDate time.Time `json:"release_date"`
	GameGenre   string    `json:"game_genre"`
}
