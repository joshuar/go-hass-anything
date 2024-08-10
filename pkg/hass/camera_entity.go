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
	*entity
	Encoding      string `json:"encoding,omitempty"`
	ImageEncoding string `json:"image_encoding,omitempty"`
	Topic         string `json:"topic"`
}

// AsCamera converts the given entity into a CameraEntity. The min, max parameters
// do not need to be specified (default min: 0, default max: 255).
func AsCamera(entity *entity) *CameraEntity {
	entity.EntityType = Camera
	entity.setTopics()
	entity.validate()

	return &CameraEntity{
		entity: entity,
		Topic:  entity.getTopicPrefix() + "/camera",
	}
}

// WithEncoding sets the encoding of the payloads. By default, a CameraEntity
// publishes raw binary data on the topic.
func (entity *CameraEntity) WithEncoding(encoding string) *CameraEntity {
	entity.Encoding = encoding

	return entity
}

// WithImageEncoding sets the image encoding of the payloads. By default, a
// CameraEntity publishes images as raw binary data on the topic. Setting a
// special value of "b64" is used by Home Assistant to represent base64 encoding
// of images.
func (entity *CameraEntity) WithImageEncoding(encoding string) *CameraEntity {
	entity.ImageEncoding = encoding

	return entity
}

func (entity *CameraEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if cfg, err = json.Marshal(entity); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(entity.ConfigTopic, cfg), nil
}
