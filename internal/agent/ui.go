// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//nolint:misspell
//revive:disable:unused-receiver
package agent

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"

	"github.com/joshuar/go-hass-anything/v9/pkg/preferences"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

var ErrInvalidPrefs = errors.New("invalid preferences")

// Preferences represents preferences the user will need to set.
type Preferences interface {
	GetValue(key string) (any, bool)
	SetValue(key string, value any) error
	GetDescription(key string) string
	IsSecret(key string) bool
	Keys() []string
}

type model struct {
	title      string
	inputs     []textinput.Model
	keys       []string
	focusIndex int
	cursorMode cursor.Mode
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

//nolint:cyclop
//revive:disable:modifies-value-receiver
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}

			cmds := make([]tea.Cmd, len(m.inputs))

			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}

			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			keyPress := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if keyPress == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if keyPress == "up" || keyPress == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))

			for idx := 0; idx <= len(m.inputs)-1; idx++ {
				if idx == m.focusIndex {
					cmds[idx] = m.inputs[idx].Focus()
					m.inputs[idx].PromptStyle = focusedStyle
					m.inputs[idx].TextStyle = focusedStyle

					continue
				}

				m.inputs[idx].Blur()
				m.inputs[idx].PromptStyle = noStyle
				m.inputs[idx].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var formOutput strings.Builder

	formOutput.WriteRune('\n')
	formOutput.WriteString(fmt.Sprintf("Set preferences for %s:\n", m.title))
	formOutput.WriteRune('\n')
	formOutput.WriteRune('\n')

	for i := range m.inputs {
		formOutput.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			formOutput.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&formOutput, "\n\n%s\n\n", *button)

	formOutput.WriteString(helpStyle.Render("cursor mode is "))
	formOutput.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	formOutput.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return formOutput.String()
}

//nolint:exhaustruct
func newPreferencesForm(name string, prefs Preferences) *model {
	model := &model{title: name, keys: prefs.Keys()}

	model.inputs = make([]textinput.Model, len(model.keys))

	for idx := range model.inputs {
		text := textinput.New()
		rawValue, found := prefs.GetValue(model.keys[idx])

		if found {
			switch value := rawValue.(type) {
			case string:
				text.SetValue(value)
			case *preferences.Preference:
				pref, ok := value.Value.(string)
				if ok {
					text.SetValue(pref)
				}
			}
		}

		text.Cursor.Style = cursorStyle
		text.CharLimit = 32
		text.PromptStyle = focusedStyle
		text.Prompt = model.keys[idx] + " > "
		text.TextStyle = focusedStyle
		model.inputs[idx] = text
	}

	return model
}

func ShowPreferences(name string, prefs Preferences) error {
	// fmt.Fprintln(os.Stdout, "data: %T", data)

	appModel := newPreferencesForm(name, prefs)

	if _, err := tea.NewProgram(appModel).Run(); err != nil {
		return fmt.Errorf("could not load preferences: %w", err)
	}

	for i := 0; i <= len(appModel.inputs)-1; i++ {
		if err := prefs.SetValue(appModel.keys[i], appModel.inputs[i].Value()); err != nil {
			log.Warn().Err(err).Str("app", name).Str("preference", appModel.keys[i]).Msg("Could not save app preference.")
		}
	}

	return nil
}
