package support

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"noname-realtime-support-chat/pkg/errors"
	"strings"
	"time"
)

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "is required"
	}
	return ""
}

func Validate(dto interface{}) error {
	validate := validator.New()

	if err := validate.Struct(dto); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.WithMessage(ErrInvalidRequest, err.Error())
		}

		var out []string
		for _, err := range err.(validator.ValidationErrors) {
			out = append(out, fmt.Sprintf("%v - %v", err.Field(), msgForTag(err.Tag())))
		}
		return errors.WithMessage(ErrInvalidRequest, strings.Join(out, ", "))
	}
	return nil
}

type DTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Status   bool   `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
