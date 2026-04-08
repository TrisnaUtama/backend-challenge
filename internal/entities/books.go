package entities

import "time"

type Books struct {
	ID        string
	Title     string
	Author    string
	Year      int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
