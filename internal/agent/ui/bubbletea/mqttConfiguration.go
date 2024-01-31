// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

var explanationText = []string{
	`Format scheme://host:port. "scheme" is one of "tcp", or "ssl", "host" is the ip-address (or hostname)
	and "port" is the port on which the broker is accepting connections.`,
	`Optional username required for authentication to the broker.`,
	`Optional password required for authentication to the broker.
	Will not be echoed to screen`,
}

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
			t.Placeholder = "scheme://host:port"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.Prompt = "Server URI (required) > "
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "User (optional) > "
			t.Placeholder = "username"
			t.CharLimit = 64
		case 2:
			t.Prompt = "Password (optional) > "
			t.Placeholder = "supersecret"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.formInputs[i] = t
	}

	return m
}
