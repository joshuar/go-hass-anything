// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joshuar/go-hass-anything/pkg/config"
	"github.com/rs/zerolog/log"
)

type Agent interface {
	AppVersion() string
	AppName() string
	AppID() string
	Stop()
	GetConfig(string, interface{}) error
	SetConfig(string, interface{}) error
}

type bubbleteaUI struct {
	p *tea.Program
}

func NewBubbleTeaUI() *bubbleteaUI {
	return &bubbleteaUI{}
}

func (ui *bubbleteaUI) ShowConfiguration(a Agent) {
	var server, username, password string

	mqttForm := mqttConfiguration()

	getValue := func(c string, v string, i int) {
		if err := a.GetConfig(c, &v); err != nil {
			log.Debug().Err(err).Msgf("Could not fetch a value for %s from agent configuration.", c)
		} else {
			mqttForm.formInputs[i].SetValue(v)
		}
	}
	getValue(config.PrefMQTTServer, server, 0)
	getValue(config.PrefMQTTUser, username, 1)
	getValue(config.PrefMQTTPassword, password, 2)

	ui.p = tea.NewProgram(mqttForm)
	if _, err := ui.p.Run(); err != nil {
		log.Error().Err(err).Msg("Could not start configuration UI.")
	}

	setValue := func(c string, i int) {
		if err := a.SetConfig(c, mqttForm.formInputs[i].Value()); err != nil {
			log.Debug().Err(err).Msgf("Could not set a value for %s in agent configuration.", c)
		}
	}
	setValue(config.PrefMQTTServer, 0)
	if mqttForm.formInputs[1].Value() != "" {
		setValue(config.PrefMQTTUser, 1)
	}
	if mqttForm.formInputs[2].Value() != "" {
		setValue(config.PrefMQTTUser, 2)
	}
}

func (ui *bubbleteaUI) Run() {
}
