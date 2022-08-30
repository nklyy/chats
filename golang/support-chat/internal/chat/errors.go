package chat

import (
	"support-chat/pkg/codes"
	"support-chat/pkg/errors"
)

const (
	StatusRequiredToken errors.Status = "token_required"
)

var (
	ErrRequiredToken = errors.New(codes.Unauthorized, StatusRequiredToken)
)
