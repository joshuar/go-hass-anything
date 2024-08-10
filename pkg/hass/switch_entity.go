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

// SwitchEntity represents an entity that can be turned on or off. For more
// details see https://www.home-assistant.io/integrations/switch.mqtt/
type SwitchEntity struct {
	*entity
	Optimistic bool `json:"optimistic,omitempty"`
}

// AsSwitch converts the given entity into a SwitchEntity. Additional builders
// can potentially be applied to customise it further.
func AsSwitch(entity *entity, optimistic bool) *SwitchEntity {
	entity.EntityType = Switch
	entity.setTopics()
	entity.validate()

	return &SwitchEntity{
		entity:     entity,
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

func (e *SwitchEntity) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}
