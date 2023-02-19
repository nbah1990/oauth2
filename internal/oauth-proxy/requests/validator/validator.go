package validator

import (
	pgValidator "github.com/go-playground/validator"
)

type PgValidatorI interface {
	Validate(i interface{}) *RError
}

type PgValidator struct {
	validator *pgValidator.Validate
}

type IError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

type RError struct {
	Errors  []IError `json:"errors"`
	Message string   `json:"message"`
}

func (pv PgValidator) Validate(i interface{}) *RError {
	err := pv.validator.Struct(i)

	if err != nil {
		var errors []IError

		for _, err := range err.(pgValidator.ValidationErrors) {
			var el IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			errors = append(errors, el)
		}

		return &RError{Errors: errors, Message: err.Error()}
	}

	return nil
}

func NewValidator() (ev *PgValidator) {
	return &PgValidator{validator: pgValidator.New()}
}
