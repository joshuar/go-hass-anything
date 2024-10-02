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
	*EntityDetails
	*EntityAttributes
	*EntityEncoding
	ImageTopic  string `json:"image_topic,omitempty" validate:"required_without=URLTopic"`
	ContentType string `json:"content_type,omitempty"`
	URLTopic    string `json:"url_topic,omitempty" validate:"required_without=ImageTopic"`
	URLTemplate string `json:"url_template,omitempty"`
	mode        ImageMode
}

// WithMode defines how the image entity operates:
//
// For ModeImage, every time a message under the image_topic in the
// configuration is received, the image displayed in Home Assistant will also be
// updated. Messages received on image_topic should contain the full contents of
// an image file, for example, a JPEG image, without any additional encoding or
// metadata.
//
// For ModeURL, the image topic defines an image URL for a new picture to show.
// The URL can be extracted from the payload by setting a template with
// WithURLTemplate.
func (e *ImageEntity) WithMode(mode ImageMode) *ImageEntity {
	e.mode = mode

	return e
}

// WithContentType defines what kind of image format the message body is using.
// For example, `image/png` or `image/jpeg`. The default is `image/jpeg`.
func (e *ImageEntity) WithContentType(contentType string) *ImageEntity {
	// ContentType is irrelevant for modes other than Image.
	if e.mode == ModeImage {
		e.ContentType = contentType
	}

	return e
}

// WithURLTemplate defines a template to use to extract the URL to the image
// from the response received. See
// https://www.home-assistant.io/docs/configuration/templating/#using-templates-with-the-mqtt-integration
// for template configuration.
func (e *ImageEntity) WithURLTemplate(template string) *ImageEntity {
	// URLTemplate is irrelevant for modes other than URL.
	if e.mode == ModeURL {
		e.URLTemplate = template
	}

	return e
}

func (e *ImageEntity) WithDetails(options ...DetailsOption) *ImageEntity {
	e.EntityDetails = WithDetails(Image, options...)

	return e
}

func (e *ImageEntity) WithEncoding(options ...EncodingOption) *ImageEntity {
	e.EntityEncoding = WithEncodingOptions(options...)

	return e
}

func (e *ImageEntity) WithAttributes(options ...AttributeOption) *ImageEntity {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *ImageEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	imageTopic := generateTopic("image", e.EntityDetails)

	switch e.mode {
	case ModeImage:
		e.ImageTopic = imageTopic
	case ModeURL:
		e.URLTopic = imageTopic
	}

	if err = validateEntity(e); err != nil {
		return nil, fmt.Errorf("entity config is invalid: %w", err)
	}

	configTopic := generateTopic("config", e.EntityDetails)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(configTopic, cfg), nil
}

func (e *ImageEntity) GetImageTopic() string {
	switch e.mode {
	case ModeImage:
		return e.ImageTopic
	case ModeURL:
		return e.URLTopic
	}

	return ""
}

func NewImageEntity() *ImageEntity {
	return &ImageEntity{}
}
