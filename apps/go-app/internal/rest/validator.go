package rest

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// FormatValidationError formats the validation error into a readable string
func FormatValidationError(err error) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		if len(ve) > 0 {
			fe := ve[0]
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				return field + " tidak boleh kosong"
			case "min":
				return field + " minimal " + fe.Param() + " karakter"
			case "max":
				return field + " maksimal " + fe.Param() + " karakter"
			case "email":
				return field + " harus berupa email yang valid"
			default:
				return field + " tidak valid"
			}
		}
	}
	return err.Error()
}
