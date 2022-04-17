package auth

import (
	"noname-support-chat/pkg/codes"
	"noname-support-chat/pkg/errors"
)

const (
	StatusInvalidRequest errors.Status = "invalid_request"
)

var (
	ErrInvalidRequest = errors.New(codes.BadRequest, StatusInvalidRequest)
)
