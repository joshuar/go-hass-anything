// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package switchapp demonstrates an app that displays a switch in Home
// Assistant.
//
//revive:disable:unused-receiver
package switchapp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/eclipse/paho.golang/paho"

	mqtthass "github.com/joshuar/go-hass-anything/v11/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

const (
	appName = "Go Hass Anything Switch Example App"
	appID   = "mqtt_switch_example"
)

type SwitchApp struct {
	entity      *mqtthass.SwitchEntity
	msgCh       chan *mqttapi.Msg
	entityState bool
}

func New(_ context.Context) (*SwitchApp, error) {
	app := &SwitchApp{
		msgCh: make(chan *mqttapi.Msg),
	}

	app.entity = mqtthass.AsSwitch(mqtthass.NewEntity(appName, "Switch", "").
		WithDeviceInfo(newDevice()).
		WithDefaultOriginInfo().
		WithCommandCallback(app.switchCommandCallback).
		WithStateCallback(app.switchStateCallback).
		WithValueTemplate("{{ value }}"),
		true).AsTypeSwitch()

	return app, nil
}

func (a *SwitchApp) Name() string {
	return appName
}

func (a *SwitchApp) Configuration() []*mqttapi.Msg {
	switchEntityCfg, err := a.entity.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal switch entity config.", "error", err)

		return nil
	}

	return []*mqttapi.Msg{switchEntityCfg}
}

func (a *SwitchApp) States() []*mqttapi.Msg {
	switchEntityState, err := a.entity.MarshalState()
	if err != nil {
		slog.Warn("Unable to marshal switch state to MQTT message.", "error", err.Error())
	}

	return []*mqttapi.Msg{switchEntityState}
}

func (a *SwitchApp) Subscriptions() []*mqttapi.Subscription {
	switchEntitySub, err := a.entity.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal switch state subscription.", "error", err.Error())
	}

	return []*mqttapi.Subscription{switchEntitySub}
}

// Update is unused, there is no app data to update.
func (a *SwitchApp) Update(_ context.Context) error { return nil }

// MsgCh returns a channel through which we could pass any message to MQTT on
// any kind of custom event trigger or other non time-based polling.
func (a *SwitchApp) MsgCh() chan *mqttapi.Msg {
	return a.msgCh
}

// switchCallback is our callback function for when the switch entity is
// manipulated in Home Assistant.
func (a *SwitchApp) switchStateCallback(_ ...any) (json.RawMessage, error) {
	switch a.entityState {
	case true:
		return json.RawMessage(`ON`), nil
	default:
		return json.RawMessage(`OFF`), nil
	}
}

// switchCommandCallback is our callback function for when the switch entity is
// manipulated in Home Assistant.
func (a *SwitchApp) switchCommandCallback(p *paho.Publish) {
	// Record the new state.
	state := string(p.Payload)
	switch state {
	case "ON":
		slog.Info("Switch was turned on.")

		a.entityState = true
	case "OFF":
		slog.Info("Switch was turned off.")

		a.entityState = false
	}
	// Publish a message with the new state.
	msg, err := a.entity.MarshalState()
	if err != nil {
		slog.Warn("Unable to marshal new state message.", "error", err.Error())
	} else {
		a.msgCh <- msg
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
