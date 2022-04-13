package user

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusUserAlreadyExists   errors.Status = "user_already_exists"
	StatusUserNotFound        errors.Status = "user_not_found"
	StatusInvalidIpAddress    errors.Status = "invalid_ip_address"
	StatusInvalidSalt         errors.Status = "invalid_salt"
	StatusFailedCreateUser    errors.Status = "failed_create_user"
	StatusFailedSaveUser      errors.Status = "failed_save_user"
	StatusFailedUpdateUser    errors.Status = "failed_update_user"
	StatusFailedFindFreeUsers errors.Status = "failed_find_free_users"
	StatusNoUsersYet          errors.Status = "no_users_yet"
)

var (
	ErrAlreadyExists       = errors.New(codes.DuplicateError, StatusUserAlreadyExists)
	ErrNotFound            = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidIpAddress    = errors.New(codes.BadRequest, StatusInvalidIpAddress)
	ErrInvalidSalt         = errors.New(codes.BadRequest, StatusInvalidSalt)
	ErrFailedCreateUser    = errors.New(codes.BadRequest, StatusFailedCreateUser)
	ErrFailedSaveUser      = errors.New(codes.BadRequest, StatusFailedSaveUser)
	ErrFailedUpdateUser    = errors.New(codes.BadRequest, StatusFailedUpdateUser)
	ErrFailedFindFreeUsers = errors.New(codes.BadRequest, StatusFailedFindFreeUsers)
	ErrNoUsersYet          = errors.New(codes.BadRequest, StatusNoUsersYet)
)
