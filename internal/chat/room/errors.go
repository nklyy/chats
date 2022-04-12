package room

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusToken             errors.Status = "invalid_token"
	StatusRequiredToken     errors.Status = "token_required"
	StatusInvalidId         errors.Status = "invalid_id"
	StatusInvalidConnection errors.Status = "invalid_connection"
	StatusInvalidName       errors.Status = "invalid_name"
	StatusInvalidUserId     errors.Status = "invalid_userId"
	StatusRoomAlreadyExists errors.Status = "room_already_exists"
	StatusUserNotFound      errors.Status = "user_not_found"
	StatusFailedCreateRoom  errors.Status = "failed_create_room"
	StatusFailedSaveRoom    errors.Status = "failed_save_room"
	StatusFailedUpdateRoom  errors.Status = "failed_update_room"
	StatusFailedDeleteRoom  errors.Status = "failed_delete_room"
)

var (
	ErrToken             = errors.New(codes.Unauthorized, StatusToken)
	ErrRequiredToken     = errors.New(codes.Unauthorized, StatusRequiredToken)
	ErrInvalidId         = errors.New(codes.BadRequest, StatusInvalidId)
	ErrInvalidConnection = errors.New(codes.BadRequest, StatusInvalidConnection)
	ErrInvalidName       = errors.New(codes.BadRequest, StatusInvalidName)
	ErrInvalidUserId     = errors.New(codes.BadRequest, StatusInvalidUserId)
	ErrAlreadyExists     = errors.New(codes.DuplicateError, StatusRoomAlreadyExists)
	ErrNotFound          = errors.New(codes.NotFound, StatusUserNotFound)
	ErrFailedCreateRoom  = errors.New(codes.BadRequest, StatusFailedCreateRoom)
	ErrFailedSaveRoom    = errors.New(codes.BadRequest, StatusFailedSaveRoom)
	ErrFailedUpdateRoom  = errors.New(codes.BadRequest, StatusFailedUpdateRoom)
	ErrFailedDeleteRoom  = errors.New(codes.BadRequest, StatusFailedDeleteRoom)
)
