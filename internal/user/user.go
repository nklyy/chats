package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	IpAddress string             `bson:"ip_address"`
	RoomName  *string            `bson:"roomName"`
	Free      bool               `bson:"free"`
	Banned    bool               `bson:"banned"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewUser(ipAddr string) (*User, error) {
	if ipAddr == "" {
		return nil, errors.WithMessage(ErrInvalidIpAddress, "should be not empty")
	}

	return &User{
		ID:        primitive.NewObjectID(),
		IpAddress: ipAddr,
		RoomName:  nil,
		Free:      true,
		Banned:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *User) SetFreeStatus(status bool) {
	s.Free = status
	s.UpdatedAt = time.Now()
}

func (s *User) SetRoom(roomName *string) {
	s.RoomName = roomName
	s.UpdatedAt = time.Now()
}

func (s *User) SetBannedStatus(status bool) {
	s.Banned = status
	s.UpdatedAt = time.Now()
}
