package user

import (
	"time"
)

type DTO struct {
	ID       string  `json:"id"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Password string  `json:"password,omitempty"`
	Support  bool    `json:"support,omitempty"`
	RoomName *string `bson:"roomName"`
	Free     bool    `bson:"free"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
