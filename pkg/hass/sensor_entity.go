// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"fmt"
	"time"

	mqttapi "github.com/joshuar/go-hass-anything/v12/pkg/mqtt"
)

// SensorEntity represents an entity which has some kind of value. For more
// details, see https://www.home-assistant.io/integrations/sensor.mqtt/
type SensorEntity struct {
	*EntityAttributes
	*EntityState
	*EntityDetails
	LastResetValueTemplate string `json:"last_reset_value_template,omitempty"`
	entityType             EntityType
	StateExpiry            int  `json:"expire_after,omitempty" validate:"omitempty,gte=0"`
	ForceUpdate            bool `json:"force_update,omitempty"`
}

// WithStateExpiry defines the number of seconds after the sensor’s state expires,
// if it’s not updated. After expiry, the sensor’s state becomes "unavailable".
func (e *SensorEntity) WithStateExpiry(expiry time.Duration) *SensorEntity {
	e.StateExpiry = int(expiry.Seconds())

	return e
}

// ForcedUpdates sends update events even if the value hasn’t changed. Useful if
// you want to have meaningful value graphs in history.
func (e *SensorEntity) ForcedUpdates() *SensorEntity {
	e.ForceUpdate = true

	return e
}

// WithLastResetValueTemplate defines a template to extract the last_reset. When
// last_reset_value_template is set, the state_class option must be total.
// Available variables: entity_id. The entity_id can be used to reference the
// entity’s attributes.
func (e *SensorEntity) WithLastResetValueTemplate(template string) *SensorEntity {
	e.LastResetValueTemplate = template

	return e
}

func (e *SensorEntity) WithDetails(options ...DetailsOption) *SensorEntity {
	e.EntityDetails = WithDetails(e.entityType, options...)

	return e
}

func (e *SensorEntity) WithState(options ...StateOption) *SensorEntity {
	e.EntityState = WithStateOptions(options...)
	e.StateTopic = generateTopic("state", e.EntityDetails)

	return e
}

func (e *SensorEntity) WithAttributes(options ...AttributeOption) *SensorEntity {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *SensorEntity) MarshalConfig() (*mqttapi.Msg, error) {
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

func NewSensorEntity() *SensorEntity {
	entity := &SensorEntity{}
	entity.entityType = Sensor

	return entity
}

func NewBinarySensorEntity() *SensorEntity {
	entity := &SensorEntity{}
	entity.entityType = BinarySensor

	return entity
}
