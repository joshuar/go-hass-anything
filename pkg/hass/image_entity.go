// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"fmt"
	"log/slog"

	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

const (
	// ModeImage is an image entity which publishes an image file.
	ModeImage ImageMode = iota
	// ModeURL is an image entity which publishes a URL to an image file.
	ModeURL
)

// ImageMode reflects how this image entity is handled by Home Assistant.
type ImageMode int

// ImageEntity represents an entity which sends image files through MQTT. For
// more details, see https://www.home-assistant.io/integrations/image.mqtt/
type ImageEntity struct {
	*entity
	Encoding      string `json:"encoding,omitempty"`
	ImageEncoding string `json:"image_encoding,omitempty"`
	ImageTopic    string `json:"image_topic,omitempty"`
	ContentType   string `json:"content_type,omitempty"`
	URLTopic      string `json:"url_topic,omitempty"`
	URLTemplate   string `json:"url_template,omitempty"`
	mode          ImageMode
}

// AsImage converts the given entity into a ImageEntity. A mode can be specified
// that controls how the entity will be handled by Home Assistant. Mode "Image"
// expects an image file as the contents of the entity's state topic. Mode "URL"
// expects a URL. The default value of mode is "Image".
func AsImage(entity *entity, mode ImageMode) *ImageEntity {
	entity.EntityType = Image
	entity.setTopics()
	entity.validate()

	imgEntity := &ImageEntity{
		entity: entity,
		mode:   mode,
	}

	switch imgEntity.mode {
	case ModeImage:
		imgEntity.ImageTopic = entity.getTopicPrefix() + "/image"
	case ModeURL:
		imgEntity.URLTopic = entity.getTopicPrefix() + "/image"
	}

	return imgEntity
}

// WithEncoding sets the encoding of the payloads. By default, a ImageEntity
// publishes raw binary data on the topic.
func (entity *ImageEntity) WithEncoding(encoding string) *ImageEntity {
	entity.Encoding = encoding

	return entity
}

// WithImageEncoding sets the image encoding of the payloads. By default, a
// ImageEntity publishes images as raw binary data on the topic. Setting a
// special value of "b64" is used by Home Assistant to represent base64 encoding
// of images.
func (entity *ImageEntity) WithImageEncoding(encoding string) *ImageEntity {
	entity.ImageEncoding = encoding

	return entity
}

// WithContentType defines what kind of image format the message body is using.
// For example, "image/png" or "image/jpeg".
func (entity *ImageEntity) WithContentType(ctype string) *ImageEntity {
	if entity.mode == ModeURL {
		slog.Warn("Ignoring content type for URL image.")

		return entity
	}

	entity.ContentType = ctype

	return entity
}

// WithURLTemplate defines a template to use to extract the URL to the image
// from the response received. See
// https://www.home-assistant.io/docs/configuration/templating/#using-templates-with-the-mqtt-integration
// for template configuration.
func (entity *ImageEntity) WithURLTemplate(template string) *ImageEntity {
	if entity.mode == ModeImage {
		slog.Warn("Ignoring url template for image.")

		return entity
	}

	entity.URLTemplate = template

	return entity
}

// GetImageTopic returns the topic on which images will appear on MQTT.
func (entity *ImageEntity) GetImageTopic() string {
	switch entity.mode {
	case ModeImage:
		return entity.ImageTopic
	case ModeURL:
		return entity.URLTopic
	}

	return ""
}

func (entity *ImageEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if cfg, err = json.Marshal(entity); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(entity.ConfigTopic, cfg), nil
}
