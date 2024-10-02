// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

//revive:disable:unchecked-type-assertion
//nolint:errorlint
func validateEntity[T any](object T) error {
	var errs error

	err := validate.Struct(object)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()
			switch {
			case err.Tag() == "required":
				errs = errors.Join(errs, fmt.Errorf("%s is required", err.Field()))
			case err.Tag() == "required_without":
				errs = errors.Join(errs, fmt.Errorf("%s cannot be set when %s is set", err.Field(), err.Param()))
			case err.StructField() == "Icon":
				errs = errors.Join(errs, errors.New("icon should be of the form 'mdi:someicon'"))
			default:
				errs = errors.Join(errs, fmt.Errorf("%s failed to validate on %s tag", err.Field(), err.Tag()))
			}
		}
	}

	return errs
}
