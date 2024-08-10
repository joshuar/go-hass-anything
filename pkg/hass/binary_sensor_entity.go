// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

// BinarySensorEntity represents an entity which has a boolean state. For more
// details, see https://www.home-assistant.io/integrations/binary_sensor.mqtt/
type BinarySensorEntity struct {
	*entity
}

// AsBinarySensor converts the given entity into a BinarySensorEntity.
// Additional builders can potentially be applied to customise it further.
func AsBinarySensor(entity *entity) *BinarySensorEntity {
	entity.EntityType = BinarySensor
	entity.setTopics()
	entity.validate()

	return &BinarySensorEntity{
		entity: entity,
	}
}
