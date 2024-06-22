// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
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

// Preferences represents preferences the user will need to set.
type Preferences interface {
	// Name is a name for this group of preferences. It could be an app name.
	Name() string
	// Preferences returns the current preferences of the app as a map[string]any
	GetPreferences() *preferences.Preferences
	// SetPreferences will set the given preferences for the app. It returns a
	// non-nil error if there was a problem setting any preferences.
	SetPreferences(prefs *preferences.Preferences) error
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
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
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
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
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
	var b strings.Builder

	b.WriteRune('\n')
	b.WriteString(fmt.Sprintf("Set preferences for %s:\n", m.title))
	b.WriteRune('\n')
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func newPreferencesForm(title string, prefs *preferences.Preferences) *model {
	model := &model{title: title, keys: prefs.Keys()}

	model.inputs = make([]textinput.Model, len(model.keys))

	for i := range model.inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Placeholder = prefs.GetString(model.keys[i])
		t.PromptStyle = focusedStyle
		t.Prompt = model.keys[i] + " > "
		t.TextStyle = focusedStyle
		model.inputs[i] = t
	}

	return model
}

func ShowPreferences(app Preferences) error {
	appPrefs := app.GetPreferences()
	appModel := newPreferencesForm(app.Name(), appPrefs)

	if _, err := tea.NewProgram(appModel).Run(); err != nil {
		return fmt.Errorf("could not load preferences: %w", err)
	}

	for i := 0; i <= len(appModel.inputs)-1; i++ {
		if err := appPrefs.Set(appModel.keys[i], appModel.inputs[i].Value()); err != nil {
			log.Warn().Err(err).Str("app", app.Name()).Str("preference", appModel.keys[i]).Msg("Could not save app preference.")
		}

		if err := app.SetPreferences(appPrefs); err != nil {
			return fmt.Errorf("could not save preferences: %w", err)
		}
	}

	return nil
}
