package auth

import (
	"support-chat/pkg/codes"
	"support-chat/pkg/errors"
)

const (
	StatusInvalidRequest errors.Status = "invalid_request"
)

var (
	ErrInvalidRequest = errors.New(codes.BadRequest, StatusInvalidRequest)
)
