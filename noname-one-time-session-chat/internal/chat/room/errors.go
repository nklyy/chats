package room

import (
	"noname-realtime-support-chat/pkg/codes"
	"noname-realtime-support-chat/pkg/errors"
)

const (
	StatusInvalidFingerprint errors.Status = "invalid_fingerprint"
	StatusInvalidConnection  errors.Status = "invalid_connection"
	StatusInvalidName        errors.Status = "invalid_name"
	StatusFailedCreateRoom   errors.Status = "failed_create_room"
)

var (
	ErrInvalidName        = errors.New(codes.BadRequest, StatusInvalidName)
	ErrFailedCreateRoom   = errors.New(codes.BadRequest, StatusFailedCreateRoom)
	ErrInvalidFingerprint = errors.New(codes.BadRequest, StatusInvalidFingerprint)
	ErrInvalidConnection  = errors.New(codes.BadRequest, StatusInvalidConnection)
)
