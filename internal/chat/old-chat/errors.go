package old_chat

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusRoomAlreadyExists errors.Status = "room_already_exists"
	StatusUserNotFound      errors.Status = "user_not_found"
	StatusFailedSaveRoom    errors.Status = "failed_save_room"
)

var (
	ErrAlreadyExists  = errors.New(codes.DuplicateError, StatusRoomAlreadyExists)
	ErrNotFound       = errors.New(codes.NotFound, StatusUserNotFound)
	ErrFailedSaveRoom = errors.New(codes.BadRequest, StatusFailedSaveRoom)
)
