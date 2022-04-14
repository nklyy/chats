package chat

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusRequiredToken     errors.Status = "token_required"
	StatusInvalidId         errors.Status = "invalid_id"
	StatusInvalidConnection errors.Status = "invalid_connection"
)

var (
	ErrRequiredToken     = errors.New(codes.Unauthorized, StatusRequiredToken)
	ErrInvalidId         = errors.New(codes.BadRequest, StatusInvalidId)
	ErrInvalidConnection = errors.New(codes.BadRequest, StatusInvalidConnection)
)
