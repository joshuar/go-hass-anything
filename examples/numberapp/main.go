// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package numberapp demonstrates an app that displays a number slider in Home
// Assistant.
//
//revive:disable:unused-receiver
package numberapp

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/eclipse/paho.golang/paho"

	mqtthass "github.com/joshuar/go-hass-anything/v12/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v12/pkg/mqtt"
)

const (
	appName = "Go Hass Anything Number Example App"
	appID   = "mqtt_number_example"

	minNum = 0
	maxNum = 250
)

type NumberApp struct {
	entity      *mqtthass.NumberEntity[int]
	entityState int
}

func New(_ context.Context) (*NumberApp, error) {
	app := &NumberApp{}

	// Our default/starting value of the entity.
	app.entityState = 50
	// Our number entity. The value can be an int or float type. The min/max can
	// be any value in the range of the type in use. We configure the entity to
	// appear as a slider in Home Assistant.
	app.entity = mqtthass.NewNumberEntity[int]().
		WithMin(minNum).WithMax(maxNum).WithStep(1).WithMode(mqtthass.NumberSlider).
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Number"),
			mqtthass.ID("number"),
			mqtthass.DeviceInfo(newDevice()),
		).
		WithState(
			mqtthass.StateCallback(app.numberStateCallback),
			mqtthass.ValueTemplate("{{ value_json.value }}"),
		).
		WithCommand(
			mqtthass.CommandCallback(app.numberCommandCallback),
		)

	return app, nil
}

func (a *NumberApp) Name() string {
	return appName
}

func (a *NumberApp) Configuration() []*mqttapi.Msg {
	numberEntityCfg, err := a.entity.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal number entity config.", "error", err)

		return nil
	}

	return []*mqttapi.Msg{numberEntityCfg}
}

func (a *NumberApp) States() []*mqttapi.Msg {
	numberEntityState, err := a.entity.MarshalState()
	if err != nil {
		slog.Warn("Unable to marshal number state to MQTT message.", "error", err.Error())
	}

	return []*mqttapi.Msg{numberEntityState}
}

func (a *NumberApp) Subscriptions() []*mqttapi.Subscription {
	numberEntitySub, err := a.entity.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal number state subscription.", "error", err.Error())
	}

	return []*mqttapi.Subscription{numberEntitySub}
}

// Update is unused, there is no app data to update.
func (a *NumberApp) Update(_ context.Context) error { return nil }

// numberStateCallback is called on the polling interval when we need to publish
// the current number state.
func (a *NumberApp) numberStateCallback(_ ...any) (json.RawMessage, error) {
	return json.RawMessage(`{ "value": ` + strconv.Itoa(a.entityState) + ` }`), nil
}

// numberCommandCallback is our callback function for when a request to change the value is
// received on MQTT, we set our state internally and publish back on the
// state topic for any listeners.
func (a *NumberApp) numberCommandCallback(p *paho.Publish) {
	if newValue, err := strconv.Atoi(string(p.Payload)); err != nil {
		slog.Warn("Could not parse new value for number.", "error", err.Error())
	} else {
		slog.Info("Number was changed.", slog.Int("value", newValue))
		a.entityState = newValue
	}
}

func newDevice() *mqtthass.Device {
	return &mqtthass.Device{
		Name:         appName,
		Identifiers:  []string{appID},
		URL:          "https://github.com/joshuar/go-hass-anything",
		Manufacturer: "go-hass-anything",
		Model:        appID,
	}
}
