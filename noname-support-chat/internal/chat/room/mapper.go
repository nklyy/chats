package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-support-chat/pkg/errors"
)

func MapToDTO(r *Model) *DTO {

	var messages []*RoomMessage
	if r.Messages != nil {
		messages = *r.Messages
	}

	return &DTO{
		ID:       r.ID.Hex(),
		Name:     r.Name,
		Messages: &messages,
	}
}

func MapToEntity(dto *DTO) (*Model, error) {
	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	return &Model{
		ID:       id,
		Name:     dto.Name,
		Messages: dto.Messages,
	}, nil
}
