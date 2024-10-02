// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:generate stringer -type=ButtonType -output button_entity_generated.go -linecomment
package hass

import (
	"encoding/json"
	"fmt"

	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

const (
	ButtonTypeNone     ButtonType = iota //
	ButtonTypeIdentify                   // identify
	ButtonTypeRestart                    // restart
	ButtonTypeUpdate                     // update
)

type ButtonType int

// ButtonEntity represents an entity which can perform some action or event in
// response to being "pushed". For more details, see
// https://www.home-assistant.io/integrations/button.mqtt/
type ButtonEntity struct {
	*EntityDetails
	*EntityCommand
	*EntityAttributes
	ButtonType   string `json:"device_class,omitempty"`
	PayloadPress string `json:"payload_press,omitempty"`
}

// WithPressPayload sets the payload to send to trigger the button. Defaults to
// PRESS.
func (e *ButtonEntity) WithPressPayload(payload string) *ButtonEntity {
	e.PayloadPress = payload

	return e
}

// WithButtonType sets the type of button, defining how it gets displayed in
// Home Assistant. See also:
// https://www.home-assistant.io/integrations/button/#device-class
func (e *ButtonEntity) WithButtonType(buttonType ButtonType) *ButtonEntity {
	e.ButtonType = buttonType.String()

	return e
}

func (e *ButtonEntity) WithDetails(options ...DetailsOption) *ButtonEntity {
	e.EntityDetails = WithDetails(Button, options...)

	return e
}

func (e *ButtonEntity) WithCommand(options ...CommandOption) *ButtonEntity {
	e.EntityCommand = WithCommandOptions(options...)
	e.CommandTopic = generateTopic("press", e.EntityDetails)

	return e
}

func (e *ButtonEntity) WithAttributes(options ...AttributeOption) *ButtonEntity {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *ButtonEntity) MarshalConfig() (*mqttapi.Msg, error) {
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

func NewButtonEntity() *ButtonEntity {
	return &ButtonEntity{}
}
