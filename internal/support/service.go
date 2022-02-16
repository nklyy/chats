package support

import (
	"errors"
	"go.uber.org/zap"
)

type Service interface {
}

type service struct {
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) (Service, error) {
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	return &service{logger: logger}, nil
}
