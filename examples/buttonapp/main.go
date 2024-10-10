// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package buttonapp demonstrates an app that displays a button in Home
// Assistant. When the button is pressed, the Home Assistant homepage will be
// opened on the device running the Go Hass Anything agent. It uses the xdg-open
// command-line tool.
//
//revive:disable:unused-receiver
package buttonapp

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/eclipse/paho.golang/paho"

	mqtthass "github.com/joshuar/go-hass-anything/v12/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v12/pkg/mqtt"
)

const (
	appName = "Go Hass Anything Button Example App"
	appID   = "mqtt_button_example"
)

type ButtonApp struct {
	entity *mqtthass.ButtonEntity
}

func New(_ context.Context) (*ButtonApp, error) {
	app := &ButtonApp{}

	app.entity = mqtthass.NewButtonEntity().
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Open HA Homepage"),
			mqtthass.ID("open_ha_homepage"),
			mqtthass.DeviceInfo(newDevice()),
		).
		WithCommand(mqtthass.CommandCallback(app.buttonCommandCallback))

	return app, nil
}

func (a *ButtonApp) Name() string {
	return appName
}

func (a *ButtonApp) Configuration() []*mqttapi.Msg {
	buttonEntityCfg, err := a.entity.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal sensor entity config.", "error", err)

		return nil
	}

	return []*mqttapi.Msg{buttonEntityCfg}
}

// States is unused. A button has no state.
func (a *ButtonApp) States() []*mqttapi.Msg { return nil }

func (a *ButtonApp) Subscriptions() []*mqttapi.Subscription {
	buttonEntitySub, err := a.entity.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal text state subscription.", "error", err.Error())
	}

	return []*mqttapi.Subscription{buttonEntitySub}
}

// Update is unused, there is no app data to update.
func (a *ButtonApp) Update(_ context.Context) error { return nil }

// buttonCommandCallback is our callback function that is run when somebody presses the
// button we created in Home Assistant. In this case, we use xdg-open to open
// the Home Assistant homepage. But we could do other things, or many things at
// once and/or react based on the response data we got (the MQTT.Message
// parameter).
func (a *ButtonApp) buttonCommandCallback(_ *paho.Publish) {
	slog.Info("Button was pressed. Opening the Home Assistant homepage.")

	if err := exec.Command("xdg-open", "https://home-assistant.io").Run(); err != nil {
		slog.Warn("Could not execute xdg-open.", "error", err.Error())
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
