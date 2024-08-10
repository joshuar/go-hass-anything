// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/constraints"

	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

//go:generate stringer -type=NumberMode -output number_entity_generated.go -linecomment
const (
	// NumberAuto will tell Home Assistant to automatically select how the number is displayed.
	NumberAuto NumberMode = 0 // auto
	// NumberBox will tell Home Assistant to display this number as a box.
	NumberBox NumberMode = 1 // box
	// NumberSlider will tell Home Assistant to display this number as a slider.
	NumberSlider NumberMode = 2 // slider
)

// NumberMode reflects how this number entity is displayed in Home Assistant. It
// can be either automatically chosen or explicitly set to display as a slider
// or box.
type NumberMode int

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

// AsNumber converts the given entity into a NumberEntity. Additional builders
// can potentially be applied to customise it further.
func AsNumber[T constraints.Ordered](entity *entity, step, min, max T, mode NumberMode) *NumberEntity[T] {
	entity.EntityType = Number
	entity.setTopics()
	entity.validate()

	return &NumberEntity[T]{
		entity: entity,
		Step:   step,
		Min:    min,
		Max:    max,
		Mode:   mode.String(),
	}
}

func (e *NumberEntity[T]) MarshalConfig() (*mqttapi.Msg, error) {
	var (
		cfg []byte
		err error
	)

	if cfg, err = json.Marshal(e); err != nil {
		return nil, fmt.Errorf("marshal config: %w", err)
	}

	return mqttapi.NewMsg(e.ConfigTopic, cfg), nil
}
