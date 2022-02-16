package support

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
	"time"
)

type Support struct {
	ID       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email"`
	Name     string             `bson:"name"`
	Password string             `bson:"password"`
	Status   bool               `bson:"status"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewSupport(email, name, password string) (*Support, error) {
	if email == "" {
		return nil, errors.WithMessage(ErrInvalidEmail, "should be not empty")
	}
	if name == "" {
		return nil, errors.WithMessage(ErrInvalidName, "should be not empty")
	}
	if password == "" {
		return nil, errors.WithMessage(ErrInvalidPassword, "should be not empty")
	}

	return &Support{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Name:      name,
		Password:  password,
		Status:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *Support) SetName(name string) {
	s.Name = name
	s.UpdatedAt = time.Now()
}

func (s *Support) SetPassword(password string) {
	s.Password = password
	s.UpdatedAt = time.Now()
}

func (s *Support) SetOnline() {
	s.Status = true
	s.UpdatedAt = time.Now()
}

func (s *Support) SetOffline() {
	s.Status = false
	s.UpdatedAt = time.Now()
}
