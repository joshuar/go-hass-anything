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
	*entity
	Mode string `json:"mode,omitempty"`
	Min  int    `json:"min,omitempty"`
	Max  int    `json:"max,omitempty"`
}

// AsText converts the given entity into a TextEntity. The min, max parameters
// do not need to be specified (default min: 0, default max: 255).
func AsText(entity *entity, min, max int) *TextEntity {
	entity.EntityType = Text
	entity.setTopics()
	entity.validate()

	if max == 0 || max > 255 {
		max = 255
	}

	return &TextEntity{
		entity: entity,
		Min:    min,
		Max:    max,
		Mode:   PlainText.String(),
	}
}

// AsPlainText sets the mode for this text entity to (the default) plain text.
func (e *TextEntity) AsPlainText() *TextEntity {
	e.Mode = PlainText.String()

	return e
}

// AsPassword sets the mode for this text entity to a password.
func (e *TextEntity) AsPassword() *TextEntity {
	e.Mode = Password.String()

	return e
}

func (e *TextEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}
