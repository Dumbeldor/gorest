package gorest

import (
	"gopkg.in/go-playground/validator.v9"
)

var gValidate *validator.Validate

// GetValidate return validator
func GetValidate() *validator.Validate {
	if gValidate == nil {
		gValidate = validator.New()
	}
	return gValidate
}
