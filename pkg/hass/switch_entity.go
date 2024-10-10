// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"fmt"

	mqttapi "github.com/joshuar/go-hass-anything/v12/pkg/mqtt"
)

// SwitchEntity represents an entity that can be turned on or off. For more
// details see https://www.home-assistant.io/integrations/switch.mqtt/
type SwitchEntity struct {
	*EntityDetails
	*EntityCommand
	*EntityState
	*EntityAttributes
	PayloadOn  string `json:"payload_on,omitempty"`
	PayloadOff string `json:"payload_off,omitempty"`
	StateOn    string `json:"state_on,omitempty"`
	StateOff   string `json:"state_off,omitempty"`
	Optimistic bool   `json:"optimistic,omitempty"`
}

// OptimisticMode ensures the switch works in optimistic mode.
func (e *SwitchEntity) OptimisticMode() *SwitchEntity {
	e.Optimistic = true

	return e
}

// WithOnPayload sets the payload that represents on state. If specified, will
// be used for both comparing to the value in the state_topic (see
// value_template and state_on for details) and sending as on command to the
// command_topic. Defaults to ON.
func (e *SwitchEntity) WithOnPayload(payload string) *SwitchEntity {
	e.PayloadOn = payload

	return e
}

// WithOffPayload sets the payload that represents off state. If specified, will
// be used for both comparing to the value in the state_topic (see
// value_template and state_off for details) and sending as off command to the
// command_topic. Defaults to OFF.
func (e *SwitchEntity) WithOffPayload(payload string) *SwitchEntity {
	e.PayloadOff = payload

	return e
}

// WithStateOn sets the payload that represents the on state. Used when value
// that represents on state in the state_topic is different from value that
// should be sent to the command_topic to turn the device on.
//
//	Default: payload_on if defined, else ON
func (e *SwitchEntity) WithStateOn(payload string) *SwitchEntity {
	e.StateOn = payload

	return e
}

// WithStateOff sets the payload that represents the off state. Used when value
// that represents off state in the state_topic is different from value that
// should be sent to the command_topic to turn the device off.
//
//	Default: payload_off if defined, else OFF
func (e *SwitchEntity) WithStateOff(payload string) *SwitchEntity {
	e.StateOff = payload

	return e
}

func (e *SwitchEntity) WithDetails(options ...DetailsOption) *SwitchEntity {
	e.EntityDetails = WithDetails(Switch, options...)

	return e
}

func (e *SwitchEntity) WithState(options ...StateOption) *SwitchEntity {
	e.EntityState = WithStateOptions(options...)
	e.StateTopic = generateTopic("state", e.EntityDetails)

	return e
}

func (e *SwitchEntity) WithCommand(options ...CommandOption) *SwitchEntity {
	e.EntityCommand = WithCommandOptions(options...)
	e.CommandTopic = generateTopic("set", e.EntityDetails)

	return e
}

func (e *SwitchEntity) WithAttributes(options ...AttributeOption) *SwitchEntity {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *SwitchEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if err = validateEntity(e); err != nil {
		return nil, fmt.Errorf("entity config is invalid: %w", err)
	}

	configTopic := generateTopic("config", e.EntityDetails)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(configTopic, cfg), nil
}

func NewSwitchEntity() *SwitchEntity {
	return &SwitchEntity{}
}
