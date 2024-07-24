// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// revive:disable:max-public-structs
// revive:disable:unexported-return

package hass

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/eclipse/paho.golang/paho"
	"github.com/iancoleman/strcase"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

var (
	ErrNoStateCallback   = errors.New("no state callback function")
	ErrNoCommandCallback = errors.New("no command callback function")
)

// HomeAssistantTopic is the prefix applied to all entity topics by default.
// Typically, this defaults to "homeassistant". It is exposed by this package
// such that it can be overridden as necessary.
var HomeAssistantTopic = "homeassistant"

// entityConfig contains fields for defining the configuration of the entity.
type entityConfig struct {
	CommandCallback    func(p *paho.Publish)
	StateCallback      func(args ...any) (json.RawMessage, error)
	AttributesCallback func(args ...any) (json.RawMessage, error)
	App                string
	NodeID             string
	EntityType         EntityType
}

// entity represents a generic entity in Home Assistant. The fields are common
// across any specific entity.
type entity struct {
	*entityConfig      `json:"-"`
	Origin             *Origin `json:"origin,omitempty"`
	Device             *Device `json:"device,omitempty"`
	UnitOfMeasurement  string  `json:"unit_of_measurement,omitempty"`
	StateClass         string  `json:"state_class,omitempty"`
	ConfigTopic        string  `json:"-"`
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
	StateExpiry        int     `json:"expire_after,omitempty"`
}

// MarshalState will generate an *mqtt.Msg for a given entity, that can be used
// to publish the entity's state to the MQTT bus.
func (e *entity) MarshalState(args ...any) (*mqttapi.Msg, error) {
	var (
		state json.RawMessage
		err   error
	)

	if e.StateCallback == nil {
		return nil, fmt.Errorf("could not marshal state for entity %s: %w", e.Name, ErrNoStateCallback)
	}

	if state, err = e.StateCallback(args...); err != nil {
		return nil, err
	}

	return mqttapi.NewMsg(e.StateTopic, state), nil
}

// MarshalAttributes will generate an *mqtt.Msg for the attributes of an entity,
// that can be used for updating the entity's attributes.
func (e *entity) MarshalAttributes(args ...any) (*mqttapi.Msg, error) {
	var (
		state json.RawMessage
		err   error
	)

	if e.AttributesCallback == nil {
		return nil, fmt.Errorf("could not marshal state for entity %s: %w", e.Name, ErrNoStateCallback)
	}

	if state, err = e.AttributesCallback(args...); err != nil {
		return nil, err
	}

	return mqttapi.NewMsg(e.AttributesTopic, state), nil
}

// MarshallSubscription will generate an *mqtt.Subscription for a given entity,
// which can be used to subscribe to an entity's command topic and execute a
// callback on messages.
func (e *entity) MarshalSubscription() (*mqttapi.Subscription, error) {
	if e.CommandCallback == nil {
		return nil, fmt.Errorf("could not marshal state for entity %s: %w", e.Name, ErrNoCommandCallback)
	}

	msg := &mqttapi.Subscription{
		Topic:    e.CommandTopic,
		Callback: e.CommandCallback,
	}

	return msg, nil
}

func (e *entity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}

// SensorEntity represents an entity which has some kind of value. For more
// details, see https://www.home-assistant.io/integrations/sensor.mqtt/
type SensorEntity struct {
	*entity
}

// BinarySensorEntity represents an entity which has a boolean state. For more
// details, see https://www.home-assistant.io/integrations/binary_sensor.mqtt/
type BinarySensorEntity struct {
	*entity
}

// ButtonEntity represents an entity which can perform some action or event in
// response to being "pushed". For more details, see
// https://www.home-assistant.io/integrations/button.mqtt/
type ButtonEntity struct {
	*entity
}

// NumberEntity represents an entity that is a number that has a given range of
// values and can be set to any value in that range, with a precision by the
// given step. For more details, see
// https://www.home-assistant.io/integrations/number.mqtt/
type NumberEntity[T constraints.Ordered] struct {
	*entity
	Min  T      `json:"min,omitempty"`
	Max  T      `json:"max,omitempty"`
	Step T      `json:"step,omitempty"`
	Mode string `json:"mode,omitempty"`
}

// SwitchEntity represents an entity that can be turned on or off. For more
// details see https://www.home-assistant.io/integrations/switch.mqtt/
type SwitchEntity struct {
	*entity
	Optimistic bool `json:"optimistic,omitempty"`
}

// TextEntity represents an entity that can display a string of text and set the
// string remotely. For more details see
// https://www.home-assistant.io/integrations/text.mqtt/
type TextEntity struct {
	*entity
	Mode string `json:"mode,omitempty"`
	Min  int    `json:"min,omitempty"`
	Max  int    `json:"max,omitempty"`
}

type EntityConstraint[T constraints.Ordered] interface {
	~*SensorEntity | ~*BinarySensorEntity | ~*ButtonEntity | ~*NumberEntity[T] | ~*SwitchEntity
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

// Topics contains the names of important topics on the MQTT bus related to the
// entity. Apps can store these topics for later retrieval and usage (for
// example, to update state topics or listen to command topics).
type Topics struct {
	Config     string
	Command    string
	State      string
	Attributes string
}

// NewEntity creates a minimal entity based on the given name and id and
// associates it with the given app. Additional builder functions should be
// chained to fill out functionality that the entity will provide.
//
//nolint:exhaustruct
func NewEntity(app, name, id string) *entity {
	name = FormatName(name)

	if id != "" {
		id = FormatID(app + "_" + id)
	} else {
		id = FormatID(app + "_" + name)
	}

	return &entity{
		entityConfig: &entityConfig{
			App: app,
		},
		UniqueID: id,
		Name:     name,
	}
}

// WithNodeID adds an additional section to the topics of the entity in MQTT. It
// can be used to help structure various entities being provided.
func (e *entity) WithNodeID(id string) *entity {
	e.NodeID = FormatID(id)

	return e
}

func (e *entity) setTopics(entityType EntityType) {
	var prefix string
	if e.NodeID != "" {
		prefix = strings.Join([]string{HomeAssistantTopic, entityType.String(), e.NodeID, e.UniqueID}, "/")
	} else {
		prefix = strings.Join([]string{HomeAssistantTopic, entityType.String(), e.UniqueID}, "/")
	}

	e.ConfigTopic = prefix + "/config"
	e.StateTopic = prefix + "/state"

	if e.CommandCallback != nil {
		e.CommandTopic = prefix + "/set"
	}

	if e.AttributesTemplate != "" || e.AttributesCallback != nil {
		e.AttributesTopic = prefix + "/attributes"
	}

	e.EntityType = entityType
}

//nolint:exhaustruct
func (e *entity) validate() {
	if e.Origin == nil {
		slog.Warn("No origin set, using default origin for entity.", "entity", e.Name)
		e.WithDefaultOriginInfo()
	}

	if e.Device == nil {
		slog.Warn("No device set, using default device for entity.", "entity", e.Name)
		e.WithDeviceInfo(&Device{
			Name:         "Go Hass Anything Default Device",
			Identifiers:  []string{"DefaultDevice"},
			URL:          "https://github.com/joshuar/go-hass-anything",
			Manufacturer: "go-hass-anything",
		})
	}
}

// WithAttributesTemplate configures the passed in template to be used to extract the
// value of the attributes in Home Assistant.
func (e *entity) WithAttributesTemplate(t string) *entity {
	e.AttributesTemplate = t

	return e
}

// WithAttributesCallback will add the passed in function as the callback action
// to be run whenever the attributes of the entity are needed. If this callback
// is to be used, then the WithAttributesTopic() builder function should also be
// called to set-up the attributes topic.
func (e *entity) WithAttributesCallback(c func(args ...any) (json.RawMessage, error)) *entity {
	e.AttributesCallback = c

	return e
}

// WithCommandCallback will add the passed in function as the callback action to
// be run when a message is received on the command topic of the entity. It
// doesn't make sense to add this for entities that don't have a command topic,
// like regular sensors.
func (e *entity) WithCommandCallback(c func(p *paho.Publish)) *entity {
	e.CommandCallback = c

	return e
}

// WithStateCallback will add the passed in function as the callback action to
// be run whenever the state of the entity is needed. It doesn't make sense to
// add this for entities that don't report a state, like buttons. It might not
// be useful to use this where you have a single state that represents many
// entities. In such cases, it would be better to manually send the state in
// your own code.
func (e *entity) WithStateCallback(c func(args ...any) (json.RawMessage, error)) *entity {
	e.StateCallback = c

	return e
}

// WithDeviceInfo adds the passed in device info to the entity config.
func (e *entity) WithDeviceInfo(d *Device) *entity {
	e.Device = d

	return e
}

// WithOriginInfo adds the passed in origin info to the entity config.
func (e *entity) WithOriginInfo(o *Origin) *entity {
	e.Origin = o

	return e
}

// WithOriginInfo adds a pre-filled origin that references go-hass-agent
// to the entity config.
//
//nolint:exhaustruct
func (e *entity) WithDefaultOriginInfo() *entity {
	e.Origin = &Origin{
		Name: "Go Hass Anything",
		URL:  "https://github.com/joshuar/go-hass-anything",
	}

	return e
}

// WithValueTemplate configures the passed in template to be used to extract the
// value of the entity in Home Assistant.
func (e *entity) WithValueTemplate(t string) *entity {
	e.ValueTemplate = t

	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "measurement".
func (e *entity) WithStateClassMeasurement() *entity {
	e.StateClass = "measurement"

	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "total".
func (e *entity) WithStateClassTotal() *entity {
	e.StateClass = "total"

	return e
}

// WithStateClassMeasurement configures the State Class for the entity to be "total_increasing".
func (e *entity) WithStateClassTotalIncreasing() *entity {
	e.StateClass = "total_increasing"

	return e
}

// WithDeviceClass configures the Device Class for the entity.
func (e *entity) WithDeviceClass(d string) *entity {
	e.DeviceClass = d

	return e
}

// WithUnits adds a unit of measurement to the entity.
func (e *entity) WithUnits(u string) *entity {
	e.UnitOfMeasurement = u

	return e
}

// WithIcon adds an icon to the entity.
func (e *entity) WithIcon(i string) *entity {
	e.Icon = i

	return e
}

// WithStateExpiry defines the number of seconds after the sensor’s state
// expires, if it’s not updated. After expiry, the sensor’s state becomes
// "unavailable".
func (e *entity) WithStateExpiry(i int) *entity {
	e.StateExpiry = i

	return e
}

// AsDiagnostic will mark this entity as a diagnostic entity in Home Assistant.
func (e *entity) AsDiagnostic() *entity {
	e.EntityCategory = "diagnostic"

	return e
}

// GetTopics returns a Topic struct containing the topics configured for this
// entity. If an entity does not have a particular topic (due to not having some
// functionality), the topic value will be an empty string.
func (e *entity) GetTopics() *Topics {
	return &Topics{
		Config:     e.ConfigTopic,
		Command:    e.CommandTopic,
		State:      e.StateTopic,
		Attributes: e.AttributesTopic,
	}
}

// AsSensor converts the given entity into a SensorEntity. Additional builders
// can potentially be applied to customise it further.
func AsSensor(e *entity) *SensorEntity {
	e.setTopics(Sensor)
	e.validate()

	return &SensorEntity{
		entity: e,
	}
}

// AsBinarySensor converts the given entity into a BinarySensorEntity.
// Additional builders can potentially be applied to customise it further.
func AsBinarySensor(e *entity) *BinarySensorEntity {
	e.setTopics(BinarySensor)
	e.validate()

	return &BinarySensorEntity{
		entity: e,
	}
}

// AsButton converts the given entity into a ButtonEntity. Additional builders
// can potentially be applied to customise it further.
func AsButton(e *entity) *ButtonEntity {
	e.setTopics(Button)
	e.validate()

	return &ButtonEntity{
		entity: e,
	}
}

// AsNumber converts the given entity into a NumberEntity. Additional builders
// can potentially be applied to customise it further.
func AsNumber[T constraints.Ordered](e *entity, step, min, max T, mode NumberMode) *NumberEntity[T] {
	e.setTopics(Number)
	e.validate()

	return &NumberEntity[T]{
		entity: e,
		Step:   step,
		Min:    min,
		Max:    max,
		Mode:   mode.String(),
	}
}

// AsSwitch converts the given entity into a SwitchEntity. Additional builders
// can potentially be applied to customise it further.
func AsSwitch(e *entity, optimistic bool) *SwitchEntity {
	e.setTopics(Switch)
	e.validate()

	return &SwitchEntity{
		entity:     e,
		Optimistic: optimistic,
	}
}

// AsTypeSwitch sets the SwitchEntity device class as a "switch". This primarily
// affects how it will be displayed in Home Assistant.
func (e *SwitchEntity) AsTypeSwitch() *SwitchEntity {
	e.DeviceClass = "switch"

	return e
}

// AsTypeSwitch sets the SwitchEntity device class as an "outlet". This
// primarily affects how it will be displayed in Home Assistant.
func (e *SwitchEntity) AsTypeOutlet() *SwitchEntity {
	e.DeviceClass = "outlet"

	return e
}

// AsText converts the given entity into a TextEntity. The min, max parameters
// do not need to be specified (default min: 0, default max: 255).
func AsText(entity *entity, min, max int) *TextEntity {
	entity.setTopics(Text)
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
