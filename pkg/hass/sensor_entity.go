// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

// SensorEntity represents an entity which has some kind of value. For more
// details, see https://www.home-assistant.io/integrations/sensor.mqtt/
type SensorEntity struct {
	*entity
}

// AsSensor converts the given entity into a SensorEntity. Additional builders
// can potentially be applied to customise it further.
func AsSensor(entity *entity) *SensorEntity {
	entity.EntityType = Sensor
	entity.setTopics()
	entity.validate()

	return &SensorEntity{
		entity: entity,
	}
}
