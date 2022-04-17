package user

import (
	"noname-support-chat/pkg/codes"
	"noname-support-chat/pkg/errors"
)

const (
	StatusUserAlreadyExists         errors.Status = "user_already_exists"
	StatusUserNotFound              errors.Status = "user_not_found"
	StatusInvalidEmail              errors.Status = "invalid_email"
	StatusInvalidName               errors.Status = "invalid_name"
	StatusInvalidPassword           errors.Status = "invalid_password"
	StatusInvalidSalt               errors.Status = "invalid_salt"
	StatusToken                     errors.Status = "invalid_token"
	StatusRequiredToken             errors.Status = "token_required"
	StatusTokenDoesntHavePermission errors.Status = "authorization_token_doesnt_have_permission"
	StatusFailedCreateUser          errors.Status = "failed_create_user"
	StatusFailedSaveUser            errors.Status = "failed_save_user"
	StatusFailedUpdateUser          errors.Status = "failed_update_user"
	StatusFailedFindFreeUsers       errors.Status = "failed_find_free_users"
	StatusNoUsersYet                errors.Status = "no_users_yet"
)

var (
	ErrAlreadyExists             = errors.New(codes.DuplicateError, StatusUserAlreadyExists)
	ErrNotFound                  = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail              = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidName               = errors.New(codes.BadRequest, StatusInvalidName)
	ErrInvalidPassword           = errors.New(codes.BadRequest, StatusInvalidPassword)
	ErrInvalidSalt               = errors.New(codes.BadRequest, StatusInvalidSalt)
	ErrToken                     = errors.New(codes.Unauthorized, StatusToken)
	ErrRequiredToken             = errors.New(codes.Unauthorized, StatusRequiredToken)
	ErrTokenDoesntHavePermission = errors.New(codes.Forbidden, StatusTokenDoesntHavePermission)
	ErrFailedCreateUser          = errors.New(codes.BadRequest, StatusFailedCreateUser)
	ErrFailedSaveUser            = errors.New(codes.BadRequest, StatusFailedSaveUser)
	ErrFailedUpdateUser          = errors.New(codes.BadRequest, StatusFailedUpdateUser)
	ErrFailedFindFreeUsers       = errors.New(codes.BadRequest, StatusFailedFindFreeUsers)
	ErrNoUsersYet                = errors.New(codes.BadRequest, StatusNoUsersYet)
)
