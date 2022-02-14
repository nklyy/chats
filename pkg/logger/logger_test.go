package logger

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLogger(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	appEnv := "development"

	tests := []struct {
		name   string
		appEnv string
		expect func(*testing.T, Logger, error)
	}{
		{
			name:   "should return logger",
			appEnv: appEnv,
			expect: func(t *testing.T, l Logger, err error) {
				assert.NotNil(t, l)
				assert.Nil(t, err)
			},
		},
		{
			name:   "should return env error",
			appEnv: "",
			expect: func(t *testing.T, l Logger, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, errors.New("invalid app env"))
				assert.Nil(t, l)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := NewLogger(tc.appEnv)
			tc.expect(t, svc, err)
		})
	}
}
