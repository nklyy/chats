package support

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusSupportAlreadyExists           errors.Status = "support_already_exists"
	StatusInvalidRequest                 errors.Status = "invalid_request"
	StatusInvalidBody                    errors.Status = "invalid_body"
	StatusUserNotFound                   errors.Status = "user_not_found"
	StatusInvalidEmail                   errors.Status = "invalid_email"
	StatusInvalidName                    errors.Status = "invalid_name"
	StatusInvalidPassword                errors.Status = "invalid_password"
	StatusInvalidSalt                    errors.Status = "invalid_salt"
	StatusAuthorizationToken             errors.Status = "invalid_authorization_token"
	StatusRequiredAuthorizationToken     errors.Status = "authorization_token_required"
	StatusFailedVerifyAuthorizationToken errors.Status = "failed_verify_authorization_token"
	StatusTokenDoesntHavePermission      errors.Status = "authorization_token_doesnt_have_permission"
)

var (
	ErrAlreadyExists                  = errors.New(codes.DuplicateError, StatusSupportAlreadyExists)
	ErrInvalidRequest                 = errors.New(codes.BadRequest, StatusInvalidRequest)
	ErrInvalidBody                    = errors.New(codes.BadRequest, StatusInvalidBody)
	ErrNotFound                       = errors.New(codes.NotFound, StatusUserNotFound)
	ErrInvalidEmail                   = errors.New(codes.BadRequest, StatusInvalidEmail)
	ErrInvalidName                    = errors.New(codes.BadRequest, StatusInvalidName)
	ErrInvalidPassword                = errors.New(codes.BadRequest, StatusInvalidPassword)
	ErrInvalidSalt                    = errors.New(codes.BadRequest, StatusInvalidSalt)
	ErrAuthorizationToken             = errors.New(codes.Forbidden, StatusAuthorizationToken)
	ErrRequiredAuthorizationToken     = errors.New(codes.Forbidden, StatusRequiredAuthorizationToken)
	ErrFailedVerifyAuthorizationToken = errors.New(codes.Forbidden, StatusFailedVerifyAuthorizationToken)
	ErrTokenDoesntHavePermission      = errors.New(codes.Forbidden, StatusTokenDoesntHavePermission)
)
