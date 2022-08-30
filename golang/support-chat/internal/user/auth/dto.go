package auth

import (
	"fmt"
	"strings"
	"support-chat/pkg/errors"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "is required"
	case "email":
		return "must be correct email"
	case "password":
		return "must be bigger than 6 characters, contain at lest one capital letter, one small letter and one number"
	}
	return ""
}

func Validate(dto interface{}) error {
	validate := validator.New()

	_ = validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		var (
			containsUpper bool
			containsLower bool
			containsDigit bool
			length        bool
		)

		if len(password) >= 6 {
			length = true
		}

		for _, char := range password {
			if unicode.IsUpper(char) {
				containsUpper = true
			} else if unicode.IsLower(char) {
				containsLower = true
			} else if unicode.IsDigit(char) {
				containsDigit = true
			}
		}
		return containsUpper && containsLower && containsDigit && length
	})

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

type RegistrationDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,password"`
}

type RegistrationResponseDTO struct {
	Id string `json:"id"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshDTO struct {
	Token string `json:"token" validate:"required"`
}

type LogoutDTO struct {
	Token string `json:"token" validate:"required"`
}
type CheckDTO struct {
	Token string `json:"token" validate:"required"`
}

type CheckResponseDTO struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	IsRoom bool   `json:"is_room"`
}
