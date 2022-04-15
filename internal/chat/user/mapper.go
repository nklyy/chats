package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
)

func MapToDTO(u *User) *DTO {
	return &DTO{
		ID:          u.ID.Hex(),
		Fingerprint: u.Fingerprint,
		RoomName:    u.RoomName,
		Banned:      u.Banned,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func MapToEntity(dto *DTO) (*User, error) {
	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	return &User{
		ID:          id,
		Fingerprint: dto.Fingerprint,
		RoomName:    dto.RoomName,
		Banned:      dto.Banned,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}, nil
}
