// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FormatName will take a string s and format it with appropriate spacing
// between words and capitalised the first letter of each word. For example
// someString becomes Some String. The new string is then an appropriate format
// to be used as a name in Home Assistant.
func FormatName(s string) string {
	c := cases.Title(language.AmericanEnglish)

	return c.String(strcase.ToDelimited(s, ' '))
}

// FormatID will take a string s and format it as snake_case. The new string is
// then an appropriate format to be used as a unique ID in Home Assistant.
func FormatID(s string) string {
	return strcase.ToSnake(s)
}
