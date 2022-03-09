package jwt

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusTokenInvalid       errors.Status = "invalid_token"
	StatusTokenExpire        errors.Status = "token_expire"
	StatusFailedCreateCache  errors.Status = "failed_create_cache"
	StatusNotFound           errors.Status = "token_not_found"
	StatusFailedExtendToken  errors.Status = "failed_extend_token"
	StatusFailedDeleteToken  errors.Status = "failed_delete_token"
	StatusFailedCreateTokens errors.Status = "failed_create_token"
)

var (
	ErrToken              = errors.New(codes.Unauthorized, StatusTokenInvalid)
	ErrTokenExpire        = errors.New(codes.Unauthorized, StatusTokenExpire)
	ErrFailedCreateCache  = errors.New(codes.Unauthorized, StatusFailedCreateCache)
	ErrNotFound           = errors.New(codes.Unauthorized, StatusNotFound)
	ErrFailedExtendToken  = errors.New(codes.Unauthorized, StatusFailedExtendToken)
	ErrFailedDeleteToken  = errors.New(codes.Unauthorized, StatusFailedDeleteToken)
	ErrFailedCreateTokens = errors.New(codes.Unauthorized, StatusFailedCreateTokens)
)
