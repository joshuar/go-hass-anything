// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func mqttConfiguration() formInputModel {
	m := formInputModel{
		formInputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.formInputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "some.host:1883"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.Prompt = "MQTT Server (required) > "
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "MQTT User (optional) > "
			t.Placeholder = "username"
			t.CharLimit = 64
		case 2:
			t.Prompt = "MQTT Password (optional) > "
			t.Placeholder = "supersecret"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.formInputs[i] = t
	}

	return m
}
