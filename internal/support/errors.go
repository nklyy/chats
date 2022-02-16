package support

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusUserNotFound    errors.Status = "user_not_found"
	StatusInvalidEmail    errors.Status = "invalid_email"
	StatusInvalidName     errors.Status = "invalid_name"
	StatusInvalidPassword errors.Status = "invalid_password"
)

var (
	ErrNotFound        = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail    = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidName     = errors.New(codes.BadRequest, StatusInvalidName)
	ErrInvalidPassword = errors.New(codes.BadRequest, StatusInvalidPassword)
)
