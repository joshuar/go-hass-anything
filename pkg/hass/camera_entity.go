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

// CameraEntity represents an entity which sends image files through MQTT. For
// more details, see https://www.home-assistant.io/integrations/camera.mqtt/
type CameraEntity struct {
	*EntityDetails
	*EntityAttributes
	*EntityEncoding
	Topic string `json:"topic" validate:"required"`
}

func (e *CameraEntity) WithDetails(options ...DetailsOption) *CameraEntity {
	e.EntityDetails = WithDetails(Camera, options...)

	return e
}

func (e *CameraEntity) WithEncoding(options ...EncodingOption) *CameraEntity {
	e.EntityEncoding = WithEncodingOptions(options...)

	return e
}

func (e *CameraEntity) WithAttributes(options ...AttributeOption) *CameraEntity {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *CameraEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if err = validateEntity(e); err != nil {
		return nil, fmt.Errorf("entity config is invalid: %w", err)
	}

	e.Topic = generateTopic("camera", e.EntityDetails)

	configTopic := generateTopic("config", e.EntityDetails)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(configTopic, cfg), nil
}

func NewCameraEntity() *CameraEntity {
	return &CameraEntity{}
}
