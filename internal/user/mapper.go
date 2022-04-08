package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
)

func MapToDTO(u *User) *DTO {
	//var roomName string
	//if u.RoomName != nil {
	//	roomName = *u.RoomName
	//}

	return &DTO{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		Support:   u.Support,
		RoomName:  u.RoomName,
		Free:      u.Free,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
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
		RoomName:  dto.RoomName,
		Free:      dto.Free,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}, nil
}
