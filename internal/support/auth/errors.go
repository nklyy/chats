package auth

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusInvalidRequest errors.Status = "invalid_request"
)

var (
	ErrInvalidRequest = errors.New(codes.BadRequest, StatusInvalidRequest)
)
