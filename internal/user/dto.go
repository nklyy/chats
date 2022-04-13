package user

import (
	"time"
)

type DTO struct {
	ID        string  `json:"id"`
	IpAddress string  `json:"ip_address"`
	RoomName  *string `json:"roomName"`
	Free      bool    `json:"free"`
	Banned    bool    `json:"banned"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
