// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"fmt"

	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

//go:generate stringer -type=TextEntityMode -output text_entity_generated.go -linecomment
const (
	PlainText TextEntityMode = iota // text
	Password                        // password
)

type TextEntityMode int

// TextEntity represents an entity that can display a string of text and set the
// string remotely. For more details see
// https://www.home-assistant.io/integrations/text.mqtt/
type TextEntity struct {
	*EntityDetails
	*EntityCommand
	*EntityAttributes
	*EntityState
	Mode    string `json:"mode,omitempty"`
	Pattern string `json:"pattern,omitempty"`
	Min     int    `json:"min,omitempty" validate:"min=0"`
	Max     int    `json:"max,omitempty" validate:"max=255"`
}

// WithMin will set the minimum size of a text being set or received.
//
//nolint:predeclared
func (e *TextEntity) WithMin(min int) *TextEntity {
	e.Min = min

	return e
}

// WithMax will set the maximum size of a text being set or received (maximum is
// 255).
//
//nolint:predeclared
func (e *TextEntity) WithMax(max int) *TextEntity {
	e.Max = max

	return e
}

// WithMode sets the mode of the text entity. Must be either PlainText or Password.
func (e *TextEntity) WithMode(mode TextEntityMode) *TextEntity {
	e.Mode = mode.String()

	return e
}

// WithPattern sets a valid regular expression the text being set or received
// must match with.
func (e *TextEntity) WithPattern(pattern string) *TextEntity {
	e.Pattern = pattern

	return e
}

func (e *TextEntity) WithDetails(options ...DetailsOption) *TextEntity {
	e.EntityDetails = WithDetails(Text, options...)

	return e
}

func (e *TextEntity) WithState(options ...StateOption) *TextEntity {
	e.EntityState = WithStateOptions(options...)
	e.StateTopic = generateTopic("state", e.EntityDetails)

	return e
}

func (e *TextEntity) WithCommand(options ...CommandOption) *TextEntity {
	e.EntityCommand = WithCommandOptions(options...)
	e.CommandTopic = generateTopic("set", e.EntityDetails)

	return e
}

func (e *TextEntity) WithAttributes(options ...AttributeOption) *TextEntity {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *TextEntity) MarshalConfig() (*mqttapi.Msg, error) {
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

func NewTextEntity() *TextEntity {
	return &TextEntity{}
}
