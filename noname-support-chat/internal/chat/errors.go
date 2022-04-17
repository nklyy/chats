package chat

import (
	"noname-support-chat/pkg/codes"
	"noname-support-chat/pkg/errors"
)

const (
	StatusRequiredToken errors.Status = "token_required"
)

var (
	ErrRequiredToken = errors.New(codes.Unauthorized, StatusRequiredToken)
)
