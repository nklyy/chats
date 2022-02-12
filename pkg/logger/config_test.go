package logger

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLogger(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tests := []struct {
		name   string
		expect func(*testing.T, Config)
	}{
		{
			name: "should return logger config",
			expect: func(t *testing.T, l Config) {
				assert.NotNil(t, l)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewLoggerConfig()
			tc.expect(t, svc)
		})
	}
}
