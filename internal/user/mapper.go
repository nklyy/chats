package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
)

func MapToDTO(s *User) *DTO {
	return &DTO{
		ID:        s.ID.Hex(),
		Email:     s.Email,
		Name:      s.Name,
		Password:  s.Password,
		Support:   s.Support,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func MapToEntity(dto *DTO) (*User, error) {
	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	return &User{
		ID:        id,
		Email:     dto.Email,
		Name:      dto.Name,
		Password:  dto.Password,
		Support:   dto.Support,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}, nil
}
