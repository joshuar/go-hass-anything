// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package sensorapp demonstrates an app that displays sensors in Home
// Assistant. When the sensor states change, their values will change in Home
// Assistant. Go Hass Anything will poll the app to update its state on the
// interval provided. The app also provides a preference for configuration of a
// URL for a data source of one of the sensors, which the agent will present
// when run with the "configure" option. The remote data is fetched through the
// use of the web helper, reducing the need for writing code to talk to a remote
// server.
//
//revive:disable:unused-receiver
package sensorapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	mqtthass "github.com/joshuar/go-hass-anything/v11/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v11/pkg/web"

	"github.com/joshuar/go-hass-anything/v11/pkg/preferences"
)

const (
	appName = "Go Hass Anything Sensors Example App"
	appID   = "mqtt_sensors_example"

	weatherURL     = "http://wttr.in/?format=j1"
	weatherURLpref = "weatherURL"
	pollInterval   = time.Minute
	pollJitter     = 5 * time.Second
)

var ErrFetchWeatherFailed = errors.New("could not get weather data")

type SensorApp struct {
	entity      *mqtthass.SensorEntity
	prefs       preferences.AppPreferences
	weatherData map[string]any
}

func New(_ context.Context) (*SensorApp, error) {
	app := &SensorApp{}

	prefs, err := preferences.LoadApp(app)
	if err != nil {
		return nil, fmt.Errorf("could not load preferences: %w", err)
	}

	app.prefs = prefs

	app.entity = mqtthass.NewSensorEntity().
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Weather Temp"),
			mqtthass.ID("weather_temp"),
			mqtthass.DeviceInfo(newDevice()),
		).
		WithState(
			mqtthass.ValueTemplate("{{ value_json.current_condition[0].temp_C }}"),
			mqtthass.StateCallback(app.weatherStateCallback),
			mqtthass.StateClassMeasurement(),
			mqtthass.Units("°C"),
			mqtthass.DeviceClass("temperature"),
		)

	return app, nil
}

func (a *SensorApp) Name() string {
	return appName
}

func (a *SensorApp) Configuration() []*mqttapi.Msg {
	sensorEntityCfg, err := a.entity.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal sensor entity config.", "error", err)

		return nil
	}

	return []*mqttapi.Msg{sensorEntityCfg}
}

func (a *SensorApp) States() []*mqttapi.Msg {
	sensorEntityState, err := a.entity.MarshalState()
	if err != nil {
		slog.Warn("Unable to marshal sensor state to MQTT message.", "error", err.Error())
	}

	return []*mqttapi.Msg{sensorEntityState}
}

// Subscriptions is unused. Sensors do not have any interactivity requirements.
func (a *SensorApp) Subscriptions() []*mqttapi.Subscription { return nil }

// Update will fetch the remote data.
func (a *SensorApp) Update(ctx context.Context) error {
	// We fetch the weather using the web.ExecuteRequest helper. As our app
	// struct satisfies both the request and response interfaces this helper
	// requires, we can pass it in.
	if err := web.ExecuteRequest(ctx, a, a); err != nil {
		return ErrFetchWeatherFailed
	}

	return nil
}

// PollConfig defines our polling interval and jitter and instructs the agent to
// fetch our state values on these.
func (a *SensorApp) PollConfig() (interval, jitter time.Duration) {
	return pollInterval, pollJitter
}

func (a *SensorApp) DefaultPreferences() preferences.AppPreferences {
	prefs := make(preferences.AppPreferences)
	prefs[weatherURLpref] = &preferences.Preference{
		Value:       weatherURL,
		Description: "The URL for the weather service to use for fetching the weather.",
		Secret:      false,
	}

	return prefs
}

// In order to use the web.ExecuteRequest helper to fetch the weather, we need
// to pass it a type that satisfies the web.Request interface. We can do this by
// adding a URL() method that returns the URL to our weather provider, to our
// app struct.
func (a *SensorApp) URL() string {
	// We get the weather service URL from our app preferences.
	if serviceURL, ok := a.prefs[weatherURLpref].Value.(string); ok {
		return serviceURL
	}
	// If we can't get the config value, log a warning and fall back to the
	// default weather URL.
	slog.Warn("Could not retrieve weather service URL from preferences.")

	return weatherURL
}

// We also need a way to save the response of the web request, and we can do
// this by satisfying the web.Response interface through adding a UnmarshalJSON
// that will take the raw response JSON and save it into our app struct.
func (a *SensorApp) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &a.weatherData); err != nil {
		return fmt.Errorf("could not parse web response: %w", err)
	}

	return nil
}

// weatherStateCallback is called on the polling interval when we need to publish
// the weather.
func (a *SensorApp) weatherStateCallback(_ ...any) (json.RawMessage, error) {
	payload, err := json.Marshal(a.weatherData)
	if err != nil {
		return nil, fmt.Errorf("could not marshal weather data to state payload: %w", err)
	}

	return payload, nil
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
