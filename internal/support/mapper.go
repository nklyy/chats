package support

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"noname-realtime-support-chat/pkg/errors"
)

func MapToDTO(s *Support) *DTO {
	return &DTO{
		ID:        s.ID.Hex(),
		Email:     s.Email,
		Name:      s.Name,
		Password:  s.Password,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func MapToEntity(dto *DTO) (*Support, error) {
	id, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	return &Support{
		ID:        id,
		Email:     dto.Email,
		Name:      dto.Name,
		Password:  dto.Password,
		Status:    dto.Status,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}, nil
}
