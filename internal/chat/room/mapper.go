package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
)

func MapToDTO(r *Model) *DTO {
	return &DTO{
		ID:   r.ID.Hex(),
		Name: r.Name,
	}
}

func MapToEntity(dto *DTO) (*Model, error) {
	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	return &Model{
		ID:   id,
		Name: dto.Name,
	}, nil
}
