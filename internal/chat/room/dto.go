package room

import "time"

type DTO struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Messages *[]*RoomMessage `json:"messages"`
}

type FormatMessages struct {
	To      string    `json:"to,omitempty"`
	From    string    `json:"from,omitempty"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
