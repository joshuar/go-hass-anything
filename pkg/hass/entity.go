// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:generate go run golang.org/x/tools/cmd/stringer -type=EntityType -output entity_generated.go -linecomment
package hass

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eclipse/paho.golang/paho"

	mqttapi "github.com/joshuar/go-hass-anything/v12/pkg/mqtt"
)

const (
	Unknown EntityType = iota // unknown
	// An entity with some kind of value, numeric or string.
	Sensor // sensor
	// An entity with a boolean value.
	BinarySensor // binary_sensor
	// An entity that changes state when activated.
	Button // button
	// An entity that is a number (float or int) with a range of values.
	Number // number
	// An entity that changes state between ON and OFF.
	Switch // switch
	// An entity that can show/set a string of text.
	Text // text
	// Any entity that can send images.
	Camera // camera
	// Any entity that can send images.
	Image // image
)

// EntityType is an iota that represents the type of entity (i.e., Switch or Sensor).
type EntityType int

var (
	ErrNoStateCallback   = errors.New("no state callback function")
	ErrNoCommandCallback = errors.New("no command callback function")
)

// HomeAssistantTopic is the prefix applied to all entity topics by default.
// Typically, this defaults to "homeassistant". It is exposed by this package
// such that it can be overridden as necessary.
var HomeAssistantTopic = "homeassistant"

type EntityAvailability struct {
	AvailabilityTopic    string `json:"availability_topic,omitempty" validate:"required"`
	AvailabilityTemplate string `json:"availability_template,omitempty"`
	PayloadAvailable     string `json:"payload_available,omitempty"`
	PayloadNotAvailable  string `json:"payload_not_available,omitempty"`
}

// EntityAttributes are the fields that can be used for entities that have
// additional attributes.
type EntityAttributes struct {
	attributesCallback func(args ...any) (json.RawMessage, error)
	// AttributesTopic defines the MQTT topic subscribed to receive a JSON
	// dictionary payload and then set as sensor attributes. Implies force_update
	// of the current sensor state when a message is received on this topic.
	AttributesTopic    string `json:"json_attributes_topic,omitempty" validate:"required"`
	AttributesTemplate string `json:"json_attributes_template,omitempty"`
}

// MarshalAttributes will generate an *mqtt.Msg for the attributes of an entity,
// that can be used for updating the entity's attributes.
func (e *EntityAttributes) MarshalAttributes(args ...any) (*mqttapi.Msg, error) {
	var (
		state json.RawMessage
		err   error
	)

	if e.attributesCallback == nil {
		return nil, fmt.Errorf("could not marshal attributes: %w", ErrNoStateCallback)
	}

	if state, err = e.attributesCallback(args...); err != nil {
		return nil, err
	}

	return mqttapi.NewMsg(e.AttributesTopic, state), nil
}

type AttributeOption func(*EntityAttributes) *EntityAttributes

func WithAttributesOptions(options ...AttributeOption) *EntityAttributes {
	entity := &EntityAttributes{}

	for _, option := range options {
		entity = option(entity)
	}

	return entity
}

// WithAttributesTemplate configures the passed in template to be used to extract the
// value of the attributes in Home Assistant.
func AttributesTemplate(t string) AttributeOption {
	return func(e *EntityAttributes) *EntityAttributes {
		e.AttributesTemplate = t

		return e
	}
}

// WithAttributesCallback will add the passed in function as the callback action
// to be run whenever the attributes of the entity are needed. If this callback
// is to be used, then the WithAttributesTopic() builder function should also be
// called to set-up the attributes topic.
func AttributesCallback(c func(args ...any) (json.RawMessage, error)) AttributeOption {
	return func(e *EntityAttributes) *EntityAttributes {
		e.attributesCallback = c

		return e
	}
}

// EntityState represents all the fields that can be used for an entity that has
// a state.
type EntityState struct {
	stateCallback func(args ...any) (json.RawMessage, error)
	// StateTopic is the MQTT topic subscribed to receive state updates. A “None” payload resets
	// to an unknown state. An empty payload is ignored.
	StateTopic         string `json:"state_topic" validate:"required"`
	ValueTemplate      string `json:"value_template"`
	UnitOfMeasurement  string `json:"unit_of_measurement,omitempty"`
	StateClass         string `json:"state_class,omitempty"`
	DeviceClass        string `json:"device_lass,omitempty"`
	SuggestedPrecision uint   `json:"suggested_display_precision,omitempty"`
}

// MarshalState will generate an *mqtt.Msg for a given entity, that can be used
// to publish the entity's state to the MQTT bus.
func (e *EntityState) MarshalState(args ...any) (*mqttapi.Msg, error) {
	var (
		state json.RawMessage
		err   error
	)

	if e.stateCallback == nil {
		return nil, fmt.Errorf("could not marshal state: %w", ErrNoStateCallback)
	}

	if state, err = e.stateCallback(args...); err != nil {
		return nil, err
	}

	return mqttapi.NewMsg(e.StateTopic, state), nil
}

// StateOption is used to add functionality to the entity state, such as
// defining the state callback function or setting the state units.
type StateOption func(*EntityState) *EntityState

// WithStateOptions will assign all of the passed in options to the EntityState.
// It will also generate an appropriate state topic using the given entity ID
// and type.
func WithStateOptions(options ...StateOption) *EntityState {
	state := &EntityState{}

	for _, option := range options {
		state = option(state)
	}

	return state
}

// Units adds a unit of measurement to the entity. Only relevant to set for
// sensors with a numeric state.
func Units(u string) StateOption {
	return func(e *EntityState) *EntityState {
		e.UnitOfMeasurement = u

		return e
	}
}

// SuggestedPrecision defines the number of decimals which should be used in the
// sensor’s state after rounding. Only relevant to set for sensors with a
// numeric state.
func SuggestedPrecision(p uint) StateOption {
	return func(e *EntityState) *EntityState {
		e.SuggestedPrecision = p

		return e
	}
}

// StateCallback will add the passed in function as the callback action to
// be run whenever the state of the entity is needed. It might not
// be useful to use this where you have a single state that represents many
// entities. In such cases, it would be better to manually send the state in
// your own code.
func StateCallback(callback func(args ...any) (json.RawMessage, error)) StateOption {
	return func(e *EntityState) *EntityState {
		e.stateCallback = callback

		return e
	}
}

// ValueTemplate configures the passed in string to be the template to be used
// to extract the value of the entity in Home Assistant.
func ValueTemplate(t string) StateOption {
	return func(e *EntityState) *EntityState {
		e.ValueTemplate = t

		return e
	}
}

// StateClassMeasurement configures the State Class for the entity to be "measurement".
func StateClassMeasurement() StateOption {
	return func(e *EntityState) *EntityState {
		e.StateClass = "measurement"

		return e
	}
}

// StateClassTotal configures the State Class for the entity to be "total".
func StateClassTotal() StateOption {
	return func(e *EntityState) *EntityState {
		e.StateClass = "total"

		return e
	}
}

// StateClassTotalIncreasing configures the State Class for the entity to be "total_increasing".
func StateClassTotalIncreasing() StateOption {
	return func(e *EntityState) *EntityState {
		e.StateClass = "total_increasing"

		return e
	}
}

// DeviceClass configures the Device Class of the entity. Device classes are
// specific to the type of entity.
func DeviceClass(class string) StateOption {
	return func(e *EntityState) *EntityState {
		e.DeviceClass = class

		return e
	}
}

type EntityDetails struct {
	Origin     *Origin `json:"origin,omitempty"`
	Device     *Device `json:"device,omitempty"`
	UniqueID   string  `json:"unique_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	app        string
	Category   string `json:"entity_category,omitempty"`
	Icon       string `json:"icon,omitempty" validate:"omitempty,startswith=mdi:"`
	entityType EntityType
	Enabled    bool `json:"enabled_by_default"`
}

type DetailsOption func(*EntityDetails) *EntityDetails

// WithDetails will assign all of the passed in details options to the entity.
// Additionally, it will convieniently set the origin info of the entity to "Go
// Hass Anything" where no custom origin is desired.
func WithDetails(entityType EntityType, options ...DetailsOption) *EntityDetails {
	details := &EntityDetails{
		entityType: entityType,
		Enabled:    true,
	}
	for _, option := range options {
		details = option(details)
	}

	// If OriginInfo option was not passed in, set the default origin.
	if details.Origin == nil {
		details = DefaultOriginInfo()(details)
	}

	return details
}

// App assigns the passed in app name to the entity.
func App(app string) DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.app = app

		return e
	}
}

// Name assigns the passed in name to the entity.
func Name(name string) DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.Name = name

		return e
	}
}

// ID assigns the passed in name to the entity. It will format the value to be
// appropriate for an entity name.
func ID(id string) DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.UniqueID = strings.ToLower(strings.ReplaceAll(id, " ", "_"))

		return e
	}
}

// Icon assigns the passed in icon string to the entity.
func Icon(icon string) DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.Icon = icon

		return e
	}
}

// AsDiagnostic ensures that the entity will appear as a diagnostic entity in
// Home Assistant.
func AsDiagnostic() DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.Category = "diagnostic"

		return e
	}
}

// NotEnabledByDefault ensures that the entity will not be enabled by default
// when first added to Home Assistant.
func NotEnabledByDefault() DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.Enabled = false

		return e
	}
}

// WithDeviceInfo adds the passed in device info to the entity config.
func DeviceInfo(d *Device) DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.Device = d

		return e
	}
}

// WithOriginInfo adds the passed in origin info to the entity config.
func OriginInfo(o *Origin) DetailsOption {
	return func(e *EntityDetails) *EntityDetails {
		e.Origin = o

		return e
	}
}

// DefaultOriginInfo adds a pre-filled origin that references go-hass-agent
// to the entity config.
func DefaultOriginInfo() DetailsOption {
	return func(entity *EntityDetails) *EntityDetails {
		entity.Origin = &Origin{
			Name: "Go Hass Anything",
			URL:  "https://github.com/joshuar/go-hass-anything",
		}

		return entity
	}
}

type EntityCommand struct {
	commandCallback func(p *paho.Publish)
	CommandTopic    string `json:"command_topic" validate:"required"`
}

type CommandOption func(*EntityCommand) *EntityCommand

func WithCommandOptions(options ...CommandOption) *EntityCommand {
	details := &EntityCommand{}
	for _, option := range options {
		details = option(details)
	}

	return details
}

func CommandCallback(callback func(p *paho.Publish)) CommandOption {
	return func(e *EntityCommand) *EntityCommand {
		e.commandCallback = callback

		return e
	}
}

// MarshallSubscription will generate an *mqtt.Subscription for a given entity,
// which can be used to subscribe to an entity's command topic and execute a
// callback on messages.
func (e *EntityCommand) MarshalSubscription() (*mqttapi.Subscription, error) {
	if e.commandCallback == nil {
		return nil, fmt.Errorf("could not marshal subscription: %w", ErrNoCommandCallback)
	}

	msg := &mqttapi.Subscription{
		Topic:    e.CommandTopic,
		Callback: e.commandCallback,
	}

	return msg, nil
}

type EntityEncoding struct {
	Encoding      string `json:"encoding,omitempty"`
	ImageEncoding string `json:"image_encoding,omitempty"`
}

type EncodingOption func(*EntityEncoding) *EntityEncoding

func WithEncodingOptions(options ...EncodingOption) *EntityEncoding {
	details := &EntityEncoding{}
	for _, option := range options {
		details = option(details)
	}

	return details
}

// WithEncoding sets the encoding of the payloads.
func WithEncoding(encoding string) EncodingOption {
	return func(e *EntityEncoding) *EntityEncoding {
		e.Encoding = encoding

		return e
	}
}

// WithImageEncoding sets the image encoding of the payloads. By default, an
// entity publishes images as raw binary data on the topic.
func WithImageEncoding(encoding string) EncodingOption {
	return func(e *EntityEncoding) *EntityEncoding {
		e.ImageEncoding = encoding

		return e
	}
}

func WithBase64ImageEncoding() EncodingOption {
	return func(e *EntityEncoding) *EntityEncoding {
		e.ImageEncoding = "b64"

		return e
	}
}

// Device contains information about the device an entity is a part of to tie it
// into the device registry in Home Assistant.
type Device struct {
	Name          string   `json:"name"`
	Manufacturer  string   `json:"manufacturer,omitempty"`
	Model         string   `json:"model,omitempty"`
	HWVersion     string   `json:"hw_version,omitempty"`
	SWVersion     string   `json:"sw_version,omitempty"`
	URL           string   `json:"configuration_url,omitempty"`
	SuggestedArea string   `json:"suggested_area,omitempty"`
	Identifiers   []string `json:"identifiers"`
	Connections   []string `json:"connections,omitempty"`
}

// Origin contains information about the app that is responsible for the entity.
// It is used by Home Assistant for logging and display purposes.
type Origin struct {
	Name    string `json:"name"`
	Version string `json:"sw_version,omitempty"`
	URL     string `json:"support_url,omitempty"`
}

// GenerateTopic takes the given topicName and entityID and generates an
// appropriate topic based on the EntityType as per the Home Assistant topic
// naming recommendations.
func generateTopic(topicName string, details *EntityDetails) string {
	appID := strings.ToLower(strings.ReplaceAll(details.app, " ", "_"))

	return strings.Join([]string{HomeAssistantTopic, details.entityType.String(), appID, details.UniqueID, topicName}, "/")
}
