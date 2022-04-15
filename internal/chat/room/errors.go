package room

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusInvalidFingerprint errors.Status = "invalid_fingerprint"
	StatusInvalidConnection  errors.Status = "invalid_connection"
	StatusInvalidName        errors.Status = "invalid_name"
	StatusRoomAlreadyExists  errors.Status = "room_already_exists"
	StatusUserNotFound       errors.Status = "user_not_found"
	StatusFailedCreateRoom   errors.Status = "failed_create_room"
	StatusFailedSaveRoom     errors.Status = "failed_save_room"
	StatusFailedUpdateRoom   errors.Status = "failed_update_room"
	StatusFailedDeleteRoom   errors.Status = "failed_delete_room"
)

var (
	ErrInvalidName        = errors.New(codes.BadRequest, StatusInvalidName)
	ErrAlreadyExists      = errors.New(codes.DuplicateError, StatusRoomAlreadyExists)
	ErrNotFound           = errors.New(codes.NotFound, StatusUserNotFound)
	ErrFailedCreateRoom   = errors.New(codes.BadRequest, StatusFailedCreateRoom)
	ErrFailedSaveRoom     = errors.New(codes.BadRequest, StatusFailedSaveRoom)
	ErrFailedUpdateRoom   = errors.New(codes.BadRequest, StatusFailedUpdateRoom)
	ErrFailedDeleteRoom   = errors.New(codes.BadRequest, StatusFailedDeleteRoom)
	ErrInvalidFingerprint = errors.New(codes.BadRequest, StatusInvalidFingerprint)
	ErrInvalidConnection  = errors.New(codes.BadRequest, StatusInvalidConnection)
)
