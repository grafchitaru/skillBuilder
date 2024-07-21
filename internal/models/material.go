package models

import "time"

type NewMaterial struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Xp           int    `json:"xp"`
	Link         string `json:"link"`
}

type Material struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserId      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Xp          int       `json:"xp"`
	Link        string    `json:"link"`
}
