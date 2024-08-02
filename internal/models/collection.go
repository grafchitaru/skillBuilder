package models

import (
	"database/sql"
	"time"
)

type NewCollection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Collection struct {
	Id          string        `json:"id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	UserId      string        `json:"user_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	SumXp       sql.NullInt64 `json:"sum_xp"`
	Xp          sql.NullInt64 `json:"xp"`
}
