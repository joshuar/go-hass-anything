// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//revive:disable:unused-receiver
package exampleapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/load"

	mqtthass "github.com/joshuar/go-hass-anything/v9/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v9/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v9/pkg/preferences"
	"github.com/joshuar/go-hass-anything/v9/pkg/web"
)

const (
	appName        = "exampleApp"
	appID          = "example_app"
	weatherURL     = "http://wttr.in/?format=j1"
	weatherURLpref = "weatherURL"
)

type ExampleApp struct {
	loadData      *load.AvgStat
	config        *preferences.Preferences
	weatherEntity *mqtthass.SensorEntity
	buttonEntity  *mqtthass.ButtonEntity
	numberEntity  *mqtthass.NumberEntity[int]
	switchEntity  *mqtthass.SwitchEntity
	msgCh         chan *mqttapi.Msg
	weatherData   map[string]any
	loadEntities  []*mqtthass.SensorEntity
	numberState   int
	switchState   bool
}

// New sets up our example app. We make use of the preference loading/saving in
// the agent to provide a file for our app preferences at
// ~/.config/go-hass-anything/exampleApp-preferences.toml. We can store whatever
// preferences our app needs in this file by providing a map[string]any that
// maps preferences to values.
func New(_ context.Context) (*ExampleApp, error) {
	app := &ExampleApp{
		msgCh: make(chan *mqttapi.Msg),
	}

	// Load the preferences from disk.
	prefs, err := preferences.LoadApp(app.Name())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("could not load %s preferences: %w", app.Name(), err)
	}

	// If there isn't already a weather provider configured, set the default one.
	if url := prefs.GetString(weatherURLpref); url == "" {
		if err := prefs.Set(weatherURLpref, weatherURL); err != nil {
			return nil, fmt.Errorf("could not set default weather url: %w", err)
		}
		log.Info().Str("app", app.Name()).Msgf("Set default weather service to %s", weatherURL)
	}

	// Save the preferences to disk.
	if err := prefs.SaveApp(app.Name()); err != nil {
		return nil, fmt.Errorf("could not save %s preferences: %w", app.Name(), err)
	}

	app.config = prefs
	return app, nil
}

// In order to use the web.ExecuteRequest helper to fetch the weather, we need
// to pass it a type that satisfies the web.Request interface. We can do this by
// adding a URL() method that returns the URL to our weather provider, to our
// app struct.
func (a *ExampleApp) URL() string {
	// we get the weather service URL from our app config. If we can't get the
	// config value, we can't continue, so we exit with an error message.
	var serviceURL string
	if serviceURL = a.config.GetString(weatherURLpref); serviceURL == "" {
		log.Warn().Msg("Could not retrieve weather service URL from config.")
	}

	return serviceURL
}

// We also need a way to save the response of the web request, and we can do
// this by satisfying the web.Response interface through adding a UnmarshalJSON
// that will take the raw response JSON and save it into our app struct.
func (a *ExampleApp) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &a.weatherData)
}

// getLoadAvgs fetches the load averages for the system running
// go-hass-anything, using the very handy gopsutil package.
func (a *ExampleApp) getLoadAvgs(ctx context.Context) error {
	var l *load.AvgStat
	var err error
	if l, err = load.AvgWithContext(ctx); err != nil {
		return err
	}
	a.loadData = l
	return nil
}

// Our app needs to satisfy the hass.MQTTDevice interface to be able to send its
// data through the agent. The following four methods achieve that.

// Name simply returns the name of this app.
func (a *ExampleApp) Name() string {
	return appName
}

// Configuration is called when our app is first registered in Home Assistant and
// will return configuration messages for the data our app will send/receive.
//
//nolint:funlen
//revive:disable:function-length
func (a *ExampleApp) Configuration() []*mqttapi.Msg {
	var msgs []*mqttapi.Msg

	// deviceInfo is used to associate each sensor to our example app device in Home Assistant
	deviceInfo := &mqtthass.Device{
		Name:         appName,
		Identifiers:  []string{"exampleApp01"},
		URL:          "https://github.com/joshuar/go-hass-anything",
		Manufacturer: "go-hass-anything",
		Model:        appID,
	}

	// originInfo is used by Home Assistant to indicate where the device/sensors
	// came from. In this case, we fill out the details to indicate Go Hass
	// Anything.
	originInfo := &mqtthass.Origin{
		Name: "Go Hass Anything",
		URL:  "https://github.com/joshuar/go-hass-anything",
	}

	// for each of our sensors, we use the builder functions in the hass package
	// to create our entity configs. These are some syntactic sugar that make it
	// a little easier to configure the myriad options for different types of
	// entities by automatically applying some defaults that are appropriate.

	// we build a sensor containing the current temperature. We rely on
	// ValueTemplate to extract out the value we are interested in. We could
	// create more sensors and extract other values out of this response if
	// desired.
	a.weatherEntity = mqtthass.AsSensor(mqtthass.NewEntity(appName, "Weather Temp", "").
		WithDeviceInfo(deviceInfo).
		WithOriginInfo(originInfo).
		WithStateClassMeasurement().
		WithDeviceClass("temperature").
		WithUnits("°C").
		WithValueTemplate("{{ value_json.current_condition[0].temp_C }}").
		WithStateCallback(a.weatherStateCallback))
	if msg, err := a.weatherEntity.MarshalConfig(); err != nil {
		log.Error().Err(err).Str("entity", a.weatherEntity.Name).Msg("Could not marshal config for entity.")
	} else {
		msgs = append(msgs, msg)
	}

	// we have three sensors for the loadavgs
	for _, l := range []string{"1", "5", "15"} {
		loadEntity := mqtthass.AsSensor(mqtthass.NewEntity(appName, "Load Average "+l+" min", "load_avg_"+l).
			WithDeviceInfo(deviceInfo).
			WithOriginInfo(originInfo).
			WithStateClassMeasurement().
			WithValueTemplate("{{ value }}").
			WithStateCallback(a.loadStateCallback))
		if msg, err := loadEntity.MarshalConfig(); err != nil {
			log.Error().Err(err).Str("entity", loadEntity.Name).Msg("Could not marshal config for entity.")
		} else {
			msgs = append(msgs, msg)
		}
		a.loadEntities = append(a.loadEntities, loadEntity)
	}

	// we also have a button that when pressed in Home Assistant, will perform
	// an action.
	a.buttonEntity = mqtthass.AsButton(mqtthass.NewEntity(appName, "Button", "").
		WithDeviceInfo(deviceInfo).
		WithOriginInfo(originInfo).
		WithCommandCallback(buttonCommandCallback))
	if msg, err := a.buttonEntity.MarshalConfig(); err != nil {
		log.Error().Err(err).Str("entity", a.buttonEntity.Name).Msg("Could not marshal config for entity.")
	} else {
		msgs = append(msgs, msg)
	}

	// we have a number slider entity. we need to track the state of the
	// slider entity, and we set a default value.
	a.numberState = 50
	// create our number entity.
	a.numberEntity = mqtthass.AsNumber(mqtthass.NewEntity(appName, "Number", "").
		WithDeviceInfo(deviceInfo).
		WithOriginInfo(originInfo).
		WithStateClassMeasurement().
		WithCommandCallback(a.numberCommandCallback).
		WithStateCallback(a.numberStateCallback).
		WithValueTemplate("{{ value_json.value }}"),
		1, 0, 250, mqtthass.NumberSlider)
	if msg, err := a.numberEntity.MarshalConfig(); err != nil {
		log.Error().Err(err).Str("entity", a.numberEntity.Name).Msg("Could not marshal config for entity.")
	} else {
		msgs = append(msgs, msg)
	}

	// we have a switch entity.
	a.switchEntity = mqtthass.AsSwitch(mqtthass.NewEntity(appName, "Switch", "").
		WithDeviceInfo(deviceInfo).
		WithOriginInfo(originInfo).
		WithCommandCallback(a.switchCommandCallback).
		WithStateCallback(a.switchStateCallback).
		WithValueTemplate("{{ value }}"),
		true).AsTypeSwitch()
	if msg, err := a.switchEntity.MarshalConfig(); err != nil {
		log.Error().Err(err).Str("entity", a.switchEntity.Name).Msg("Could not marshal config for entity.")
	} else {
		msgs = append(msgs, msg)
	}

	// return the array of config messages, letting the receiver handle publishing.
	return msgs
}

// States is called when we want to send our sensor data to Home Assistant.
func (a *ExampleApp) States() []*mqttapi.Msg {
	var msgs []*mqttapi.Msg

	// we retrieve the weather data and create a message to publish it to its
	// state topic.
	weatherState, err := a.weatherEntity.MarshalState()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to marshal weather state to MQTT message.")
	} else {
		msgs = append(msgs, weatherState)
	}

	// we retrieve our load avgs and reate msgs to publish to each of their
	// state topics.
	for i, l := range []string{"1", "5", "15"} {
		loadState, err := a.loadEntities[i].MarshalState(l)
		if err != nil {
			log.Warn().Err(err).Msg("Unable to marshal load state to MQTT message.")
		} else {
			msgs = append(msgs, loadState)
		}
	}

	// we create a msg to publish the current number value to its state topic
	numberState, err := a.numberEntity.MarshalState()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to marshal number state to MQTT message.")
	} else {
		msgs = append(msgs, numberState)
	}

	// we create a msg to publish the current state of the switch
	switchState, err := a.switchEntity.MarshalState()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to marshal switch state to MQTT message.")
	} else {
		msgs = append(msgs, switchState)
	}

	// return all the messages for publication.
	return msgs
}

// Subscriptions is called once to register the callbacks we will use when Home
// Assistant sends back messages on any command topics.
func (a *ExampleApp) Subscriptions() []*mqttapi.Subscription {
	var msgs []*mqttapi.Subscription

	// we add a subscription to watch for requests for changing the number
	// value.
	buttonSub, err := a.buttonEntity.MarshalSubscription()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to marshal button subscription.")
	} else {
		msgs = append(msgs, buttonSub)
	}
	// we add a subscription to watch for requests for changing the number
	// value.
	numberSub, err := a.numberEntity.MarshalSubscription()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to marshal number subscription.")
	} else {
		msgs = append(msgs, numberSub)
	}
	// we add a subscription to watch for requests for changing the switch
	// value.
	switchSub, err := a.switchEntity.MarshalSubscription()
	if err != nil {
		log.Warn().Err(err).Msg("Unable to marshal switch subscription.")
	} else {
		msgs = append(msgs, switchSub)
	}

	// return all the messages for publication.
	return msgs
}

// Run is the function that the agent calls to start our app. In it, we create
// our app struct, register our app (if needed), listen for our button press,
// then set up a loop to send our sensor data.
func (a *ExampleApp) Update(ctx context.Context) error {
	var errs error
	// We fetch the weather using the web.ExecuteRequest helper. As our app
	// struct satisfies both the request and response interfaces this helper
	// requires, we can pass it in.
	if err := web.ExecuteRequest(ctx, a, a); err != nil {
		errs = errors.Join(errs, errors.New("could not get weather data"))
	}
	// get the load averages
	if err := a.getLoadAvgs(ctx); err != nil {
		errs = errors.Join(errs, errors.New("could not get load averages"))
	}
	return errs
}

func (a *ExampleApp) PollConfig() (interval, jitter time.Duration) {
	return time.Minute, time.Second * 5
}

func (a *ExampleApp) MsgCh() chan *mqttapi.Msg {
	return a.msgCh
}

func (a *ExampleApp) GetPreferences() *preferences.Preferences {
	return a.config
}

func (a *ExampleApp) SetPreferences(prefs *preferences.Preferences) error {
	return prefs.SaveApp(a.Name())
}

// weatherStateCallback is called on the polling interval when we need to publish
// the weather.
func (a *ExampleApp) weatherStateCallback(_ ...any) (json.RawMessage, error) {
	return json.Marshal(a.weatherData)
}

// loadStateCallback is called on the polling interval when we need to publish
// the current load averages.
func (a *ExampleApp) loadStateCallback(args ...any) (json.RawMessage, error) {
	var value float64
	var loadType string
	var ok bool

	if loadType, ok = args[0].(string); !ok {
		return nil, errors.New("could not determine which load was requested")
	}
	switch loadType {
	case "1":
		value = a.loadData.Load1
	case "5":
		value = a.loadData.Load5
	case "15":
		value = a.loadData.Load15
	}
	return json.RawMessage(strconv.FormatFloat(value, 'f', -1, 64)), nil
}

// buttonCommandCallback is our callback function that is run when somebody presses the
// button we created in Home Assistant. In this case, we use xdg-open to open
// the Home Assistant homepage. But we could do other things, or many things at
// once and/or react based on the response data we got (the MQTT.Message
// parameter).
func buttonCommandCallback(_ *paho.Publish) {
	log.Info().Msg("Button was pressed.")
	log.Info().Msg("Opening Home Assistant homepage.")
	if err := exec.Command("xdg-open", "https://home-assistant.io").Run(); err != nil {
		log.Warn().Err(err).Msg("Could not execute xdg-open.")
	}
}

// numberStateCallback is called on the polling interval when we need to publish
// the current number state.
func (a *ExampleApp) numberStateCallback(_ ...any) (json.RawMessage, error) {
	return json.RawMessage(`{ "value": ` + strconv.Itoa(a.numberState) + ` }`), nil
}

// numberCommandCallback is our callback function for when a request to change the value is
// received on MQTT, we set our state internally and publish back on the
// state topic for any listeners.
func (a *ExampleApp) numberCommandCallback(p *paho.Publish) {
	if newValue, err := strconv.Atoi(string(p.Payload)); err != nil {
		log.Warn().Err(err).Msg("Could not parse new value for number.")
	} else {
		log.Info().Int("value", newValue).Msg("Number was changed.")
		a.numberState = newValue
	}
}

// switchCallback is our callback function for when the switch entity is
// manipulated in Home Assistant.
func (a *ExampleApp) switchStateCallback(_ ...any) (json.RawMessage, error) {
	switch a.switchState {
	case true:
		return json.RawMessage(`ON`), nil
	default:
		return json.RawMessage(`OFF`), nil
	}
}

// switchCommandCallback is our callback function for when the switch entity is
// manipulated in Home Assistant.
func (a *ExampleApp) switchCommandCallback(p *paho.Publish) {
	// Record the new state.
	state := string(p.Payload)
	switch state {
	case "ON":
		log.Info().Msg("Switch was turned on.")
		a.switchState = true
	case "OFF":
		log.Info().Msg("Switch was turned off.")
		a.switchState = false
	}
	// Publish a message with the new state.
	msg, err := a.switchEntity.MarshalState()
	if err != nil {
		log.Warn().Msg("Unable to marshal new state message.")
	} else {
		a.msgCh <- msg
	}
}
