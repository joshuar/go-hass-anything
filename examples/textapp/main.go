// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package textapp demostrates an app that displays a text area in Home
// Assistant. Any text entered in the textarea will be shown as a notification
// on the device running the Go Hass Anything the notify-send command-line tool.
//
//revive:disable:unused-receiver
package textapp

import (
	"context"
	"encoding/json"
	"log/slog"
	"os/exec"

	"github.com/eclipse/paho.golang/paho"

	mqtthass "github.com/joshuar/go-hass-anything/v11/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

const (
	appName = "Go Hass Anything Text Area Example App"
	appID   = "mqtt_textentity_example"

	minTextLen = 0
	maxTextLen = 255
)

type TextApp struct {
	entity *mqtthass.TextEntity
	text   string
}

func New(_ context.Context) (*TextApp, error) {
	app := &TextApp{}
	app.text = "Replace this text and hit return!"
	app.entity = mqtthass.NewTextEntity().
		WithMode(mqtthass.PlainText).
		WithMin(minTextLen).
		WithMax(maxTextLen).
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Text"),
			mqtthass.ID("text"),
			mqtthass.DeviceInfo(newDevice()),
		).
		WithCommand(
			mqtthass.CommandCallback(app.commandCallback),
		).
		WithState(
			mqtthass.StateCallback(app.stateCallback),
			mqtthass.ValueTemplate("{{ value }}"),
		)

	return app, nil
}

func (a *TextApp) Name() string {
	return appName
}

func (a *TextApp) Configuration() []*mqttapi.Msg {
	textEntityCfg, err := a.entity.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal text entity config.",
			slog.Any("error", err))

		return nil
	}

	return []*mqttapi.Msg{textEntityCfg}
}

func (a *TextApp) States() []*mqttapi.Msg {
	textEntityState, err := a.entity.MarshalState()
	if err != nil {
		slog.Warn("Unable to marshal text state to MQTT message.",
			slog.Any("error", err))
	}

	return []*mqttapi.Msg{textEntityState}
}

func (a *TextApp) Subscriptions() []*mqttapi.Subscription {
	textEntitySub, err := a.entity.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal text state subscription.",
			slog.Any("error", err))
	}

	return []*mqttapi.Subscription{textEntitySub}
}

// Update is unused, there is no app data to update.
func (a *TextApp) Update(_ context.Context) error { return nil }

func newDevice() *mqtthass.Device {
	return &mqtthass.Device{
		Name:         appName,
		Identifiers:  []string{appID},
		URL:          "https://github.com/joshuar/go-hass-anything",
		Manufacturer: "go-hass-anything",
		Model:        appID,
	}
}

func (a *TextApp) stateCallback(_ ...any) (json.RawMessage, error) {
	return json.RawMessage(a.text), nil
}

func (a *TextApp) commandCallback(p *paho.Publish) {
	a.text = string(p.Payload)
	if err := exec.Command("notify-send", a.text).Run(); err != nil {
		slog.Warn("Could not execute notify-send.",
			slog.Any("error", err))
	}
}
