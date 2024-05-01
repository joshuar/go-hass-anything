// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"errors"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	mqttapi "github.com/joshuar/go-hass-anything/v8/pkg/mqtt"
)

// EntityConfig is a type used by apps to represent a Home Assistant entity and
// some additional fields for manipulating that entity in the app. Specify the
// type T to indicate the type of value this entity uses, which will be used to
// configure any additional fields/parameters that need to be created for the
// entity type. For example, for a number entity, specifying T as int will
// ensure the min, max and step parameters are also treated as ints.
type EntityConfig[T any] struct {
	Entity             *Entity[T]
	App                string
	CommandCallback    func(mqtt.Client, mqtt.Message)
	StateCallback      func(args ...any) (json.RawMessage, error)
	AttributesCallback func() (json.RawMessage, error)
	ConfigTopic        string
	topicPrefix        string
}

// Entity represents a generic entity in Home Assistant. The fields are common
// across any specific entity.
type Entity[T any] struct {
	numberEntity[T]
	Origin             *Origin `json:"origin,omitempty"`
	Device             *Device `json:"device,omitempty"`
	UnitOfMeasurement  string  `json:"unit_of_measurement,omitempty"`
	StateClass         string  `json:"state_class,omitempty"`
	CommandTopic       string  `json:"command_topic,omitempty"`
	StateTopic         string  `json:"state_topic"`
	ValueTemplate      string  `json:"value_template"`
	UniqueID           string  `json:"unique_id"`
	Name               string  `json:"name"`
	EntityCategory     string  `json:"entity_category,omitempty"`
	Icon               string  `json:"icon,omitempty"`
	AttributesTopic    string  `json:"json_attributes_topic,omitempty"`
	AttributesTemplate string  `json:"json_attributes_template,omitempty"`
	DeviceClass        string  `json:"device_class,omitempty"`
}

// numberEntity represents the fields specifically for number entities in Home
// Assistant.
type numberEntity[T any] struct {
	Min  T      `json:"min,omitempty"`
	Max  T      `json:"max,omitempty"`
	Step T      `json:"step,omitempty"`
	Mode string `json:"mode,omitempty"`
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
	Name    string `json:"name,omitempty"`
	Version string `json:"sw_version,omitempty"`
	URL     string `json:"support_url,omitempty"`
}

// Topics is a helper struct that is returned when an EntityConfig is created
// containing the names of important topics on the MQTT bus related to the
// entity. Apps can store these topics for later retrieval and usage (for
// example, to update state topics or listen to command topics).
type Topics struct {
	Config     string
	Command    string
	State      string
	Attributes string
}

// MarshalConfig will generate an *mqtt.Msg for a given entity, that can be used
// to publish the required config for the entity to the MQTT bus.
func MarshalConfig[T any](e *EntityConfig[T]) (*mqttapi.Msg, error) {
	var cfg []byte
	var err error
	if cfg, err = json.Marshal(e.Entity); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}

// MarshalState will generate an *mqtt.Msg for a given entity, that can be used
// to publish the entity's state to the MQTT bus.
func MarshalState[T any](e *EntityConfig[T]) (*mqttapi.Msg, error) {
	var state json.RawMessage
	var err error
	if e.StateCallback == nil {
		return nil, errors.New("entity does not have a state callback function")
	}
	if state, err = e.StateCallback(); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.Entity.StateTopic, state), nil
}

// MarshallSubscription will generate an *mqtt.Subscription for a given entity,
// which can be used to subscribe to an entity's command topic and execute a
// callback on messages.
func MarshalSubscription[T any](e *EntityConfig[T]) (*mqttapi.Subscription, error) {
	if e.CommandCallback == nil {
		return nil, errors.New("entity does not have a command callback function")
	}
	msg := &mqttapi.Subscription{
		Topic:    e.Entity.CommandTopic,
		Callback: e.CommandCallback,
	}
	return msg, nil
}

// NewEntityByName will create a new entity and config based off the given name
// and app. Use this function where you don't care about the id of the
// underlying sensor in Home Assistant. The id will be derived from the name by
// converting it to snake_case.
func NewEntityByName[T any](name, app, prefix string) *EntityConfig[T] {
	return &EntityConfig[T]{
		Entity: &Entity[T]{
			Name:     FormatName(name),
			UniqueID: FormatID(name),
		},
		App:         strings.ToLower(app),
		topicPrefix: prefix,
	}
}

// NewEntityByID will create a new entity and config based off the given id and
// app. Use this when you want to ensure the exact format of the id for the
// underlying sensor in Home Assistant. The name will be derived from the id.
func NewEntityByID[T any](id, app, prefix string) *EntityConfig[T] {
	return &EntityConfig[T]{
		Entity: &Entity[T]{
			Name:     FormatName(id),
			UniqueID: id,
		},
		App:         strings.ToLower(app),
		topicPrefix: prefix,
	}
}

// AsSensor will configure appropriate MQTT topics to represent a Home Assistant sensor.
func (e *EntityConfig[T]) AsSensor() *EntityConfig[T] {
	prefix := strings.Join([]string{e.topicPrefix, "sensor", e.App, e.Entity.UniqueID}, "/")
	e.ConfigTopic = prefix + "/config"
	e.Entity.StateTopic = prefix + "/state"
	e.Entity.AttributesTopic = prefix + "/attributes"
	return e
}

// AsBinarySensor will configure appropriate MQTT topics to represent a Home Assistant binary_sensor.
func (e *EntityConfig[T]) AsBinarySensor() *EntityConfig[T] {
	prefix := strings.Join([]string{e.topicPrefix, "binary_sensor", e.App, e.Entity.UniqueID}, "/")
	e.ConfigTopic = prefix + "/config"
	e.Entity.StateTopic = prefix + "/state"
	e.Entity.AttributesTopic = prefix + "/attributes"
	return e
}

// AsButton will configure appropriate MQTT topics to represent a Home Assistant button.
func (e *EntityConfig[T]) AsButton() *EntityConfig[T] {
	prefix := strings.Join([]string{e.topicPrefix, "button", e.App, e.Entity.UniqueID}, "/")
	e.ConfigTopic = prefix + "/config"
	e.Entity.CommandTopic = prefix + "/toggle"
	e.Entity.AttributesTopic = prefix + "/attributes"
	return e
}

// AsNumber will configure appropriate MQTT topics to represent a Home Assistant
// number. See also https://www.home-assistant.io/integrations/number.mqtt/
func (e *EntityConfig[T]) AsNumber(step, min, max T, mode NumberMode) *EntityConfig[T] {
	prefix := strings.Join([]string{e.topicPrefix, "number", e.App, e.Entity.UniqueID}, "/")
	e.ConfigTopic = prefix + "/config"
	e.Entity.CommandTopic = prefix + "/set"
	e.Entity.StateTopic = prefix + "/state"
	e.Entity.AttributesTopic = prefix + "/attributes"
	e.Entity.Step = step
	e.Entity.Min = min
	e.Entity.Max = max
	e.Entity.Mode = mode.String()
	return e
}

// WithAttributesTemplate configures the passed in template to be used to extract the
// value of the attributes in Home Assistant.
func (e *EntityConfig[T]) WithAttributesTemplate(t string) *EntityConfig[T] {
	e.Entity.AttributesTemplate = t
	return e
}

// WithAttributesCallback will add the passed in function as the callback action
// to be run whenever the attributes of the entity are needed. If this callback
// is to be used, then the WithAttributesTopic() builder function should also be
// called to set-up the attributes topic.
func (e *EntityConfig[T]) WithAttributesCallback(c func() (json.RawMessage, error)) *EntityConfig[T] {
	e.AttributesCallback = c
	return e
}

// WithCommandCallback will add the passed in function as the callback action to
// be run when a message is received on the command topic of the entity. It
// doesn't make sense to add this for entities that don't have a command topic,
// like regular sensors.
func (e *EntityConfig[T]) WithCommandCallback(c func(mqtt.Client, mqtt.Message)) *EntityConfig[T] {
	e.CommandCallback = c
	return e
}

// WithStateCallback will add the passed in function as the callback action to
// be run whenever the state of the entity is needed. It doesn't make sense to
// add this for entities that don't report a state, like buttons. It might not
// be useful to use this where you have a single state that represents many
// entities. In such cases, it would be better to manually send the state in
// your own code.
func (e *EntityConfig[T]) WithStateCallback(c func(args ...any) (json.RawMessage, error)) *EntityConfig[T] {
	e.StateCallback = c
	return e
}

// WithDeviceInfo adds the passed in device info to the entity config.
func (e *EntityConfig[T]) WithDeviceInfo(d *Device) *EntityConfig[T] {
	e.Entity.Device = d
	return e
}

// WithOriginInfo adds the passed in origin info to the entity config.
func (e *EntityConfig[T]) WithOriginInfo(o *Origin) *EntityConfig[T] {
	e.Entity.Origin = o
	return e
}

// WithOriginInfo adds a pre-filled origin that references go-hass-agent
// to the entity config.
func (e *EntityConfig[T]) WithDefaultOriginInfo() *EntityConfig[T] {
	e.Entity.Origin = &Origin{
		Name: "Go Hass Anything",
		URL:  "https://github.com/joshuar/go-hass-anything",
	}
	return e
}

// WithValueTemplate configures the passed in template to be used to extract the
// value of the entity in Home Assistant.
func (e *EntityConfig[T]) WithValueTemplate(t string) *EntityConfig[T] {
	e.Entity.ValueTemplate = t
	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "measurement".
func (e *EntityConfig[T]) WithStateClassMeasurement() *EntityConfig[T] {
	e.Entity.StateClass = "measurement"
	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "total".
func (e *EntityConfig[T]) WithStateClassTotal() *EntityConfig[T] {
	e.Entity.StateClass = "total"
	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "total_increasing".
func (e *EntityConfig[T]) WithStateClassTotalIncreasing() *EntityConfig[T] {
	e.Entity.StateClass = "total_increasing"
	return e
}

// WithDeviceClass configures the Device Class for the entity.
func (e *EntityConfig[T]) WithDeviceClass(d string) *EntityConfig[T] {
	e.Entity.DeviceClass = d
	return e
}

// WithUnits adds a unit of measurement to the entity.
func (e *EntityConfig[T]) WithUnits(u string) *EntityConfig[T] {
	e.Entity.UnitOfMeasurement = u
	return e
}

// WithIcon adds an icon to the entity.
func (e *EntityConfig[T]) WithIcon(i string) *EntityConfig[T] {
	e.Entity.Icon = i
	return e
}

// GetTopics returns a Topic struct containing the topics configured for this
// entity. If an entity does not have a particular topic (due to not having some
// functionality), the topic value will be an empty string.
func (e *EntityConfig[T]) GetTopics() *Topics {
	return &Topics{
		Config:     e.ConfigTopic,
		Command:    e.Entity.CommandTopic,
		State:      e.Entity.StateTopic,
		Attributes: e.Entity.AttributesTopic,
	}
}

// FormatName will take a string s and format it with appropriate spacing
// between words and capitalised the first letter of each word. For example
// someString becomes Some String. The new string is then an appropriate format
// to be used as a name in Home Assistant.
func FormatName(s string) string {
	c := cases.Title(language.AmericanEnglish)
	return c.String(strcase.ToDelimited(s, ' '))
}

// FormatID will take a string s and format it as snake_case. The new string is
// then an appropriate format to be used as a unique ID in Home Assistant.
func FormatID(s string) string {
	return strcase.ToSnake(s)
}
