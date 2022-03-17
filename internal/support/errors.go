package support

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusSupportAlreadyExists      errors.Status = "support_already_exists"
	StatusUserNotFound              errors.Status = "user_not_found"
	StatusInvalidEmail              errors.Status = "invalid_email"
	StatusInvalidName               errors.Status = "invalid_name"
	StatusInvalidPassword           errors.Status = "invalid_password"
	StatusInvalidSalt               errors.Status = "invalid_salt"
	StatusToken                     errors.Status = "invalid_token"
	StatusRequiredToken             errors.Status = "token_required"
	StatusTokenDoesntHavePermission errors.Status = "authorization_token_doesnt_have_permission"
	StatusFailedCreateSupport       errors.Status = "failed_create_support"
	StatusFailedSaveSupport         errors.Status = "failed_save_support"
	StatusFailedUpdateSupport       errors.Status = "failed_update_support"
)

var (
	ErrAlreadyExists             = errors.New(codes.DuplicateError, StatusSupportAlreadyExists)
	ErrNotFound                  = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail              = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidName               = errors.New(codes.BadRequest, StatusInvalidName)
	ErrInvalidPassword           = errors.New(codes.BadRequest, StatusInvalidPassword)
	ErrInvalidSalt               = errors.New(codes.BadRequest, StatusInvalidSalt)
	ErrToken                     = errors.New(codes.Unauthorized, StatusToken)
	ErrRequiredToken             = errors.New(codes.Unauthorized, StatusRequiredToken)
	ErrTokenDoesntHavePermission = errors.New(codes.Forbidden, StatusTokenDoesntHavePermission)
	ErrFailedCreateSupport       = errors.New(codes.BadRequest, StatusFailedCreateSupport)
	ErrFailedSaveSupport         = errors.New(codes.BadRequest, StatusFailedSaveSupport)
	ErrFailedUpdateSupport       = errors.New(codes.BadRequest, StatusFailedUpdateSupport)
)
