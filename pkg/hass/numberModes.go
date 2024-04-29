// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

//go:generate stringer -type=NumberMode -output numberModeStrings.go -linecomment
const (
	// NumberAuto will tell Home Assistant to automatically select how the number is displayed.
	NumberAuto NumberMode = 0 // auto
	// NumberBox will tell Home Assistant to display this number as a box.
	NumberBox NumberMode = 1 // box
	// NumberSlider will tell Home Assistant to display this number as a slider.
	NumberSlider NumberMode = 2 // slider
)

// NumberMode reflects how this number entity is displayed in Home Assistant. It
// can be either automatically chosen or explicitly set to display as a slider
// or box.
type NumberMode int
