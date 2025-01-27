package models

import "time"

type Asset struct {
	ID     int64
	UserID int64

	Type      string
	Sticker   string
	Body      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
