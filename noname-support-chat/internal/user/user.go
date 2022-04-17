package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"noname-support-chat/pkg/errors"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email"`
	Name     string             `bson:"name"`
	Password string             `bson:"password"`
	Support  bool               `bson:"support"`
	RoomName *string            `bson:"roomName"`
	Free     bool               `bson:"free"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewUser(email, name, password string, salt *int) (*User, error) {
	if email == "" {
		return nil, errors.WithMessage(ErrInvalidEmail, "should be not empty")
	}
	if name == "" {
		return nil, errors.WithMessage(ErrInvalidName, "should be not empty")
	}
	if password == "" {
		return nil, errors.WithMessage(ErrInvalidPassword, "should be not empty")
	}
	if salt == nil {
		return nil, errors.WithMessage(ErrInvalidSalt, "should be not empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), *salt)
	if err != nil {
		return nil, errors.WithMessage(ErrInvalidPassword, err.Error())
	}

	return &User{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Name:      name,
		Password:  string(hashedPassword),
		Support:   false,
		RoomName:  nil,
		Free:      true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *User) SetName(name string) {
	s.Name = name
	s.UpdatedAt = time.Now()
}

func (s *User) SetFreeStatus(status bool) {
	s.Free = status
	s.UpdatedAt = time.Now()
}

func (s *User) SetRoom(roomName *string) {
	s.RoomName = roomName
	s.UpdatedAt = time.Now()
}

func (s *User) SetPassword(password string) {
	s.Password = password
	s.UpdatedAt = time.Now()
}

func (s *User) RemovePassword() {
	s.Password = ""
}

func (s *User) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
	if err != nil {
		return false, ErrInvalidPassword
	}

	return true, nil
}
