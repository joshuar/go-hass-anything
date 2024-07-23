// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

//go:generate stringer -type=EntityType -output entityTypeStrings.go -linecomment
const (
	Unknown EntityType = iota // unknown
	// An entity with some kind of value, numeric or string.
	Sensor // sensor
	// An entity with a boolean value.
	BinarySensor // binary_sensor
	// An entity that changes state when activated.
	Button // button
	// An entity that is a number (float or int) with a range of values.
	Number // number
	// An entity that changes state between ON and OFF.
	Switch // switch
	// An entity that can show/set a string of text.
	Text // text
)

type EntityType int
