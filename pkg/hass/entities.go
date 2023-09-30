// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"errors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"github.com/joshuar/go-hass-anything/pkg/mqtt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EntityConfig struct {
	Entity          *Entity
	App             string
	CommandCallback func(MQTT.Client, MQTT.Message)
	StateCallback   func() (json.RawMessage, error)
	ConfigTopic     string
}

type Entity struct {
	Origin            *Origin `json:"origin,omitempty"`
	Device            *Device `json:"device,omitempty"`
	DeviceClass       string  `json:"device_class,omitempty"`
	StateTopic        string  `json:"state_topic"`
	StateClass        string  `json:"state_class,omitempty"`
	CommandTopic      string  `json:"command_topic,omitempty"`
	UnitOfMeasurement string  `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string  `json:"value_template"`
	UniqueID          string  `json:"unique_id"`
	Name              string  `json:"name"`
	EntityCategory    string  `json:"entity_category,omitempty"`
}

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

type Origin struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"sw_version,omitempty"`
	URL     string `json:"support_url,omitempty"`
}

// MarshalConfig will marshal a config message and payload from the given
// EntityConfig.
func MarshalConfig(e *EntityConfig) (*mqtt.MQTTMsg, error) {
	var msg *mqtt.MQTTMsg
	if jsonConfig, err := json.Marshal(e.Entity); err != nil {
		return nil, err
	} else {
		msg = &mqtt.MQTTMsg{
			Topic:    e.ConfigTopic,
			Message:  jsonConfig,
			Retained: true,
		}
	}
	return msg, nil
}

// MarshalState will marshal a state message and payload from the given
// EntityConfig and state value. Where an entity state is combined with other
// entities, it might be better to manually create a state message.
func MarshalState(e *EntityConfig) (*mqtt.MQTTMsg, error) {
	if e.StateCallback == nil {
		return nil, errors.New("entity does not have a state callback function")
	}
	if value, err := e.StateCallback(); err != nil {
		return nil, err
	} else {
		msg := &mqtt.MQTTMsg{
			Topic:    e.Entity.StateTopic,
			Message:  value,
			Retained: false,
		}
		return msg, nil
	}
}

// MarshalState will marshal a subscription from the given EntityConfig.
func MarshalSubscription(e *EntityConfig) (*mqtt.MQTTSubscription, error) {
	if e.CommandCallback == nil {
		return nil, errors.New("entity does not have a command callback function")
	}
	msg := &mqtt.MQTTSubscription{
		Topic:    e.Entity.CommandTopic,
		Callback: e.CommandCallback,
	}
	return msg, nil
}

// NewEntityByName will create a new entity and config based off the given name
// and app. Use this function where you don't care about the id of the
// underlying sensor in Home Assistant. The id will be derived from the name by
// converting it to snake_case.
func NewEntityByName(name, app string) *EntityConfig {
	return &EntityConfig{
		Entity: &Entity{
			Name:     FormatName(name),
			UniqueID: FormatID(name),
		},
		App: app,
	}
}

// NewEntityByID will create a new entity and config based off the given id and
// app. Use this when you want to ensure the exact format of the id for the
// underlying sensor in Home Assistant. The name will be derived from the id.
func NewEntityByID(id, app string) *EntityConfig {
	return &EntityConfig{
		Entity: &Entity{
			Name:     FormatName(id),
			UniqueID: id,
		},
		App: app,
	}
}

// AsSensor will configure appropriate MQTT topics to represent a Home Assistant sensor.
func (e *EntityConfig) AsSensor() *EntityConfig {
	e.ConfigTopic = mqtt.DiscoveryPrefix + "/sensor/" + e.App + "/" + e.Entity.UniqueID + "/config"
	e.Entity.StateTopic = mqtt.DiscoveryPrefix + "/sensor/" + e.App + "/" + e.Entity.UniqueID + "/state"
	return e
}

// AsBinarySensor will configure appropriate MQTT topics to represent a Home Assistant binary_sensor.
func (e *EntityConfig) AsBinarySensor() *EntityConfig {
	e.ConfigTopic = mqtt.DiscoveryPrefix + "/binary_sensor/" + e.App + "/" + e.Entity.UniqueID + "/config"
	e.Entity.StateTopic = mqtt.DiscoveryPrefix + "/binary_sensor/" + e.App + "/" + e.Entity.UniqueID + "/state"
	return e
}

// AsButton will configure appropriate MQTT topics to represent a Home Assistant button.
func (e *EntityConfig) AsButton() *EntityConfig {
	e.ConfigTopic = mqtt.DiscoveryPrefix + "/button/" + e.App + "/" + e.Entity.UniqueID + "/config"
	e.Entity.CommandTopic = mqtt.DiscoveryPrefix + "/button/" + e.App + "/" + e.Entity.UniqueID + "/toggle"
	return e
}

// WithCommandCallback will add the passed in function as the callback action to
// be run when a message is received on the command topic of the entity. It
// doesn't make sense to add this for entities that don't have a command topic,
// like regular sensors.
func (e *EntityConfig) WithCommandCallback(c func(MQTT.Client, MQTT.Message)) *EntityConfig {
	e.CommandCallback = c
	return e
}

// WithStateCallback will add the passed in function as the callback action to
// be run whenever the state of the entity is needed. It doesn't make sense to
// add this for entities that don't report a state, like buttons. It might not
// be useful to use this where you have a single state that represents many
// entities. In such cases, it would be better to manually send the state in
// your own code.
func (e *EntityConfig) WithStateCallback(c func() (json.RawMessage, error)) *EntityConfig {
	e.StateCallback = c
	return e
}

// WithDevice adds the passed in device info to the entity config.
func (e *EntityConfig) WithDevice(d *Device) *EntityConfig {
	e.Entity.Device = d
	return e
}

// WithOrigin adds the passed in origin info to the entity config.
func (e *EntityConfig) WithOrigin(o *Origin) *EntityConfig {
	e.Entity.Origin = o
	return e
}

// WithValueTemplate configures the passed in template to be used to extract the
// value of the entity in Home Assistant.
func (e *EntityConfig) WithValueTemplate(t string) *EntityConfig {
	e.Entity.ValueTemplate = t
	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "measurement".
func (e *EntityConfig) WithStateClassMeasurement() *EntityConfig {
	e.Entity.StateClass = "measurement"
	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "total".
func (e *EntityConfig) WithStateClassTotal() *EntityConfig {
	e.Entity.StateClass = "total"
	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "total_increasing".
func (e *EntityConfig) WithStateClassTotalIncreasing() *EntityConfig {
	e.Entity.StateClass = "total_increasing"
	return e
}

// WithDeviceClass configures the Device Class for the entity.
func (e *EntityConfig) WithDeviceClass(d string) *EntityConfig {
	e.Entity.DeviceClass = d
	return e
}

// WithUnits adds a unit of measurement ot the entity.
func (e *EntityConfig) WithUnits(u string) *EntityConfig {
	e.Entity.UnitOfMeasurement = u
	return e
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
