package chat

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusRequiredToken errors.Status = "token_required"
)

var (
	ErrRequiredToken = errors.New(codes.Unauthorized, StatusRequiredToken)
)
