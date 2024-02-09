// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"

	"github.com/joshuar/go-hass-anything/v5/pkg/preferences"
)

type Agent interface {
	AppVersion() string
	AppName() string
	AppID() string
	Stop()
}

type bubbleteaUI struct {
	p *tea.Program
}

func NewBubbleTeaUI() *bubbleteaUI {
	return &bubbleteaUI{}
}

func (ui *bubbleteaUI) ShowConfiguration() {
	mqttForm := mqttConfiguration()

	prefs, err := preferences.LoadPreferences()
	if err != nil {
		log.Warn().Err(err).Msg("No existing config found.")
	} else {
		mqttForm.formInputs[0].SetValue(prefs.Server)
		mqttForm.formInputs[1].SetValue(prefs.User)
		mqttForm.formInputs[2].SetValue(prefs.Password)
	}

	ui.p = tea.NewProgram(mqttForm)
	if _, err = ui.p.Run(); err != nil {
		log.Error().Err(err).Msg("Could not start configuration UI.")
	}

	err = preferences.SavePreferences(
		preferences.MQTTServer(mqttForm.formInputs[0].Value()),
		preferences.MQTTUser(mqttForm.formInputs[1].Value()),
		preferences.MQTTPassword(mqttForm.formInputs[2].Value()),
	)
	if err != nil {
		log.Error().Err(err).Msg("Problem saving preferences.")
	}
}

func (ui *bubbleteaUI) Run() {}
