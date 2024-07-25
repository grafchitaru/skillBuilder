package models

import "time"

type TypeMaterial struct {
	Id             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Characteristic string    `json:"characteristic"`
	Xp             int       `json:"xp"`
}
