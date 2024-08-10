// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

// ButtonEntity represents an entity which can perform some action or event in
// response to being "pushed". For more details, see
// https://www.home-assistant.io/integrations/button.mqtt/
type ButtonEntity struct {
	*entity
}

// AsButton converts the given entity into a ButtonEntity. Additional builders
// can potentially be applied to customise it further.
func AsButton(entity *entity) *ButtonEntity {
	entity.EntityType = Button
	entity.setTopics()
	entity.validate()

	return &ButtonEntity{
		entity: entity,
	}
}
