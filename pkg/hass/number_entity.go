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

//go:generate go run golang.org/x/tools/cmd/stringer -type=NumberMode -output number_entity_generated.go -linecomment
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
	Min  T `json:"min,omitempty"`
	Max  T `json:"max,omitempty"`
	Step T `json:"step,omitempty"`
	*EntityDetails
	*EntityState
	*EntityCommand
	*EntityAttributes
	Mode         string `json:"mode,omitempty"`
	ResetPayload string `json:"payload_reset,omitempty"`
	Optimistic   bool   `json:"optimistic,omitempty"`
}

// OptimisticMode ensures the number entity works in optimistic mode.
func (e *NumberEntity[T]) OptimisticMode() *NumberEntity[T] {
	e.Optimistic = true

	return e
}

// WithMin sets the minimum value for the number entity state.
//
//nolint:predeclared
func (e *NumberEntity[T]) WithMin(min T) *NumberEntity[T] {
	e.Min = min

	return e
}

// WithMax sets the maximum value for the number entity state.
//
//nolint:predeclared
func (e *NumberEntity[T]) WithMax(max T) *NumberEntity[T] {
	e.Max = max

	return e
}

// WithStep sets the step value. For entities with float value, smallest step
// accepted by Home Assistant 0.001.
func (e *NumberEntity[T]) WithStep(step T) *NumberEntity[T] {
	e.Step = step

	return e
}

// WithMode controls how the number should be displayed in the Home Assistant
// UI. Can be set to NumberBox or NumberSlide to force a display mode. Default
// is NumberAuto which will let Home Assistant decide.
func (e *NumberEntity[T]) WithMode(mode NumberMode) *NumberEntity[T] {
	e.Mode = mode.String()

	return e
}

// WithResetPayload defines a special payload that resets the state to unknown
// when received on the state_topic.
func (e *NumberEntity[T]) WithResetPayload(payload string) *NumberEntity[T] {
	e.ResetPayload = payload

	return e
}

func (e *NumberEntity[T]) WithDetails(options ...DetailsOption) *NumberEntity[T] {
	e.EntityDetails = WithDetails(Number, options...)

	return e
}

func (e *NumberEntity[T]) WithState(options ...StateOption) *NumberEntity[T] {
	e.EntityState = WithStateOptions(options...)
	e.StateTopic = generateTopic("state", e.EntityDetails)

	return e
}

func (e *NumberEntity[T]) WithCommand(options ...CommandOption) *NumberEntity[T] {
	e.EntityCommand = WithCommandOptions(options...)
	e.CommandTopic = generateTopic("set", e.EntityDetails)

	return e
}

func (e *NumberEntity[T]) WithAttributes(options ...AttributeOption) *NumberEntity[T] {
	e.EntityAttributes = WithAttributesOptions(options...)
	e.AttributesTopic = generateTopic("attributes", e.EntityDetails)

	return e
}

func (e *NumberEntity[T]) MarshalConfig() (*mqttapi.Msg, error) {
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

func NewNumberEntity[T constraints.Ordered]() *NumberEntity[T] {
	return &NumberEntity[T]{}
}
