package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
	"time"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id"`
	Fingerprint string             `bson:"fingerprint"`
	RoomName    *string            `bson:"room_name"`
	Banned      bool               `bson:"banned"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewUser(fingerprint, salt string) (*User, error) {
	if fingerprint == "" {
		return nil, errors.WithMessage(ErrInvalidIpAddress, "should be not empty")
	}
	if salt == "" {
		return nil, errors.WithMessage(ErrInvalidSalt, "should be not empty")
	}

	return &User{
		ID:          primitive.NewObjectID(),
		Fingerprint: fingerprint,
		RoomName:    nil,
		Banned:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

//func (s *User) SetFreeStatus(status bool) {
//	s.Free = status
//	s.UpdatedAt = time.Now()
//}

func (s *User) SetRoom(roomName *string) {
	s.RoomName = roomName
	s.UpdatedAt = time.Now()
}

func (s *User) SetBannedStatus(status bool) {
	s.Banned = status
	s.UpdatedAt = time.Now()
}
