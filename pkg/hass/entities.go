// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eclipse/paho.golang/paho"
	"github.com/iancoleman/strcase"
	"golang.org/x/exp/constraints"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	mqttapi "github.com/joshuar/go-hass-anything/v9/pkg/mqtt"
)

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
	var state json.RawMessage
	var err error
	if e.StateCallback == nil {
		return nil, fmt.Errorf("entity %s for app %s does not have a state callback function", e.UniqueID, e.App)
	}
	if state, err = e.StateCallback(args...); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.StateTopic, state), nil
}

func (e *entity) MarshalAttributes(args ...any) (*mqttapi.Msg, error) {
	var state json.RawMessage
	var err error
	if e.AttributesCallback == nil {
		return nil, fmt.Errorf("entity %s for app %s does not have a state callback function", e.UniqueID, e.App)
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
		return nil, fmt.Errorf("entity %s for app %s does not have a command callback function", e.UniqueID, e.App)
	}
	msg := &mqttapi.Subscription{
		Topic:    e.CommandTopic,
		Callback: e.CommandCallback,
	}
	return msg, nil
}

type SensorEntity struct {
	*entity
}

func (e *SensorEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var cfg []byte
	var err error
	if cfg, err = json.Marshal(e); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}

type BinarySensorEntity struct {
	*entity
}

func (e *BinarySensorEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var cfg []byte
	var err error
	if cfg, err = json.Marshal(e); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}

type ButtonEntity struct {
	*entity
}

func (e *ButtonEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var cfg []byte
	var err error
	if cfg, err = json.Marshal(e); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}

type NumberEntity[T constraints.Ordered] struct {
	*entity
	Min  T      `json:"min,omitempty"`
	Max  T      `json:"max,omitempty"`
	Step T      `json:"step,omitempty"`
	Mode string `json:"mode,omitempty"`
}

func (e *NumberEntity[T]) MarshalConfig() (*mqttapi.Msg, error) {
	var cfg []byte
	var err error
	if cfg, err = json.Marshal(e); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}

type SwitchEntity struct {
	*entity
	Optimistic bool `json:"optimistic,omitempty"`
}

func (e *SwitchEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var cfg []byte
	var err error
	if cfg, err = json.Marshal(e); err != nil {
		return nil, err
	}
	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
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

func (e *entity) WithNodeID(id string) *entity {
	e.NodeID = FormatID(id)
	return e
}

func (e *entity) setTopics(t EntityType) {
	var prefix string
	if e.NodeID != "" {
		prefix = strings.Join([]string{HomeAssistantTopic, t.String(), e.NodeID, e.UniqueID}, "/")
	} else {
		prefix = strings.Join([]string{HomeAssistantTopic, t.String(), e.UniqueID}, "/")
	}
	e.ConfigTopic = prefix + "/config"
	e.StateTopic = prefix + "/state"
	if e.CommandCallback != nil {
		e.CommandTopic = prefix + "/set"
	}
	if e.AttributesTemplate != "" || e.AttributesCallback != nil {
		e.AttributesTopic = prefix + "/attributes"
	}
	e.EntityType = t
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

func (e *entity) WithStateExpiry(i int) *entity {
	e.StateExpiry = i
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

func AsSensor(e *entity) *SensorEntity {
	e.setTopics(Sensor)
	return &SensorEntity{
		entity: e,
	}
}

func AsBinarySensor(e *entity) *BinarySensorEntity {
	e.setTopics(BinarySensor)
	return &BinarySensorEntity{
		entity: e,
	}
}

func AsButton(e *entity) *ButtonEntity {
	e.setTopics(Button)
	return &ButtonEntity{
		entity: e,
	}
}

func AsNumber[T constraints.Ordered](e *entity, step, min, max T, mode NumberMode) *NumberEntity[T] {
	e.setTopics(Number)
	return &NumberEntity[T]{
		entity: e,
		Step:   step,
		Min:    min,
		Max:    max,
		Mode:   mode.String(),
	}
}

func AsSwitch(e *entity, optimistic bool) *SwitchEntity {
	e.setTopics(Switch)
	return &SwitchEntity{
		entity:     e,
		Optimistic: optimistic,
	}
}

func (e *SwitchEntity) AsTypeSwitch() *SwitchEntity {
	e.DeviceClass = "switch"
	return e
}

func (e *SwitchEntity) AsTypeOutlet() *SwitchEntity {
	e.DeviceClass = "outlet"
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
