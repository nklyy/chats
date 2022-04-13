package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/scrypt"
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

func NewUser(ipAddr, salt string) (*User, error) {
	if ipAddr == "" {
		return nil, errors.WithMessage(ErrInvalidIpAddress, "should be not empty")
	}
	if salt == "" {
		return nil, errors.WithMessage(ErrInvalidSalt, "should be not empty")
	}

	hashAddr, err := scrypt.Key([]byte(ipAddr), []byte(salt), 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        primitive.NewObjectID(),
		IpAddress: string(hashAddr),
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
