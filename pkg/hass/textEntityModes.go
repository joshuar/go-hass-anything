// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

//go:generate stringer -type=TextEntityMode -output textEntityModeStrings.go -linecomment
const (
	PlainText TextEntityMode = iota // text
	Password                        // password
)

type TextEntityMode int
