// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package exampleapp

import (
	"context"
	"encoding/json"
	"os/exec"
	"strconv"
	"time"

	"github.com/carlmjohnson/requests"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joshuar/go-hass-anything/internal/apps/helpers"
	"github.com/joshuar/go-hass-anything/pkg/config"
	viper "github.com/joshuar/go-hass-anything/pkg/config/viper"
	"github.com/joshuar/go-hass-anything/pkg/hass"
	"github.com/joshuar/go-hass-anything/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/pkg/web"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/load"
)

const (
	appName        = "exampleApp"
	weatherURL     = "http://wttr.in/?format=j1"
	weatherURLpref = "weatherURL"
)

type exampleApp struct {
	config      config.App
	loadData    *load.AvgStat
	weatherData []byte
}

// newExampleApp sets up our example app. We make use of our custom viper
// wrapper that allows storing our config in
// ~/.config/go-hass-anything/exampleApp/config.toml. This is used for storing
// the registration status of our app (whether we need to send config messages
// or not at a minimum) and we also store the URL to the weather service.
func newExampleApp() *exampleApp {
	var err error
	app := &exampleApp{}
	// load our app config. if we don't have a config, set some defaults
	if app.config, err = viper.Load(appName); err != nil && app.config != nil {
		log.Info().Msgf("Setting default weather service to %s", weatherURL)
		app.config.Set(weatherURLpref, weatherURL)
	}
	return app
}

// getWeather sends a request for the weather, using the web.ExecuteRequest
// function. ExecuteRequest returns a chan web.Response which we can wait on for
// the web request to complete. Once it does, the web.Response has some methods
// to extract the response headers and body and any error that might have been
// returned. In our case, we don't care about the headers, just the body (and we
// check/handle any errors). We could do some processing of the response body in
// this function if needed, but in this case, it is just a JSON blob and we can
// pass that directly to MQTT and let Home Assistant extract values from it.
func (a *exampleApp) getWeather(ctx context.Context) error {
	// we make our web request using web.ExecuteRequest and wait for the
	// response, handling any error returned.
	responseCh := make(chan web.Response, 1)
	web.ExecuteRequest(ctx, a, responseCh)
	r := <-responseCh
	if r.Error() != nil {
		return r.Error()
	}

	// we save the raw JSON response as a []byte, which we can pass directly to MQTT
	// and have Home Assistant extract out the values we want.
	a.weatherData = r.Body().Bytes()
	return nil
}

// In order to use web.ExecuteRequest, we need to pass it a type that satisfies
// the web.Request interface, which has a Builder() method. Builder returns a
// *requests.Builder used to execute the actual web request. We can make our
// exampleApp struct satisfy web.Request by giving it a Builder function. In our
// case, we just need to build a request with a URL. But we could do more fancy
// requests if needed.
func (a *exampleApp) Builder() *requests.Builder {
	// we get the weather service URL from our app config. If we can't get the
	// config value, we can't continue, so we exit with an error message.
	var serviceURL string
	err := a.config.Get(weatherURLpref, &serviceURL)
	if err != nil {
		log.Error().Err(err).Msg("Could not retrieve weather service URL from config.")
		return nil
	}

	return requests.URL(serviceURL)
}

// web.Request also has a Timeout() method, which we can use to specify a
// timeout after which our web request will be cancelled. In our case, we will
// wait 15 seconds.
func (a *exampleApp) Timeout() time.Duration {
	return 15 * time.Second
}

// getLoadAvgs fetches the load averages for the system running
// go-hass-anything, using the very handy gopsutil package
func (a *exampleApp) getLoadAvgs(ctx context.Context) error {
	if l, err := load.AvgWithContext(ctx); err != nil {
		return err
	} else {
		a.loadData = l
		return nil
	}
}

// Our app needs to satisfy the hass.MQTTDevice interface to be able to send its
// data through the agent. The following three methods achieve that.

// Configuration is called when our app is first registered in Home Assistant and
// will return configuration messages for the data our app will send/receive.
func (a *exampleApp) Configuration() []*mqtt.MQTTMsg {
	var msgs []*mqtt.MQTTMsg
	var entities []*hass.EntityConfig

	// deviceInfo is used to associate each sensor to our example app device in Home Assistant
	deviceInfo := &hass.Device{
		Name:        appName,
		Identifiers: []string{"exampleApp01"},
	}

	// originInfo is used by Home Assistant to indicate where the device/sensors
	// came from. In this case, we fill out the details to indicate Go Hass
	// Anything.
	originInfo := &hass.Origin{
		Name:    appName,
		Version: "1.0.0",
		URL:     "https://github.com/joshuar/go-hass-anything",
	}

	// for each of our sensors, we use the builder functions in the hass package
	// to create our entity configs. These are some syntactic sugar that make it
	// a little easier to configure the myriad options for different types of
	// entities by automatically applying some defaults that are appropriate.

	// we build a sensor containing the current temperature. We rely on
	// ValueTemplate to extract out the value we are interested in. We could
	// create more sensors and extract other values out of this response if
	// desired.
	entities = append(entities,
		hass.NewEntityByName("ExampleApp Weather Temp", appName).
			AsSensor().
			WithDevice(deviceInfo).
			WithOrigin(originInfo).
			WithStateClassMeasurement().
			WithDeviceClass("temperature").
			WithUnits("°C").
			WithValueTemplate("{{ value_json.current_condition[0].temp_C }}"))

	// we have three sensors for the loadavgs
	for _, load := range []string{"1", "5", "15"} {
		entities = append(entities,
			hass.NewEntityByID("example_app_load"+load, appName).
				AsSensor().
				WithDevice(deviceInfo).
				WithOrigin(originInfo).
				WithStateClassMeasurement().
				WithValueTemplate("{{ value }}"))
	}

	// we also have a button that when pressed in Home Assistant, will perform
	// an action
	entities = append(entities,
		hass.NewEntityByID("example_app_button", appName).
			AsButton().
			WithCommandCallback(buttonCallback))

	// we marshal our configs into an mqtt.MQTTMsg
	for _, e := range entities {
		if msg, err := hass.MarshalConfig(e); err != nil {
			log.Error().Err(err).Str("entity", e.Entity.Name).Msg("Could not marshal config for entity.")
		} else {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

// States is called when we want to send our sensor data to Home Assistant
func (a *exampleApp) States() []*mqtt.MQTTMsg {
	var msgs []*mqtt.MQTTMsg

	// we retrieve the weather data and send that as the weather sensor state
	msgs = append(msgs, &mqtt.MQTTMsg{
		Topic:   mqtt.DiscoveryPrefix + "/sensor/" + appName + "/example_app_weather_temp/state",
		Message: a.weatherData,
	})

	// we retrieve our load avgs
	for _, load := range []string{"1", "5", "15"} {
		id := "example_app_load" + load
		var l float64
		switch load {
		case "1":
			l = a.loadData.Load1
		case "5":
			l = a.loadData.Load5
		case "15":
			l = a.loadData.Load15
		}
		msgs = append(msgs, &mqtt.MQTTMsg{
			Topic:   mqtt.DiscoveryPrefix + "/sensor/" + appName + "/" + id + "/state",
			Message: json.RawMessage(strconv.FormatFloat(l, 'f', -1, 64)),
		})
	}

	return msgs
}

// Subscriptions is called once to register the callbacks we will use when Home
// Assistant sends back messages on any command topics.
func (a *exampleApp) Subscriptions() []*mqtt.MQTTSubscription {
	var msgs []*mqtt.MQTTSubscription

	// we add our callback for our button
	msgs = append(msgs, &mqtt.MQTTSubscription{
		Topic:    mqtt.DiscoveryPrefix + "/button/" + appName + "/example_app_button/toggle",
		Callback: buttonCallback,
	})
	return msgs
}

// Run is the function that the agent calls to start our app. In it, we create
// our app struct, register our app (if needed), listen for our button press,
// then set up a loop to send our sensor data.
func Run(ctx context.Context, client hass.MQTTClient) {
	log.Info().Str("appName", appName).Msg("Starting app.")
	app := newExampleApp()

	// check if our app is registered. If not, send the configuration messages
	// and then register the app.
	if !app.config.IsRegistered(appName) {
		if err := hass.PublishConfigs(app, client); err != nil {
			log.Error().Err(err).Msg("Failed to publish configs.")
			return
		} else {
			if err = app.config.Register(appName); err != nil {
				log.Error().Err(err).Msg("Failed to register configuration.")
				return
			}
		}
	}

	// add our button subscription
	if err := hass.Subscribe(app, client); err != nil {
		log.Error().Err(err).Msg("Could not activate subscriptions.")
	}

	// create a function we will use whenever we want to send our app state
	sendState := func() {
		// get the weather
		if err := app.getWeather(ctx); err != nil {
			log.Error().Err(err).Msg("Could not get weather data.")
		}
		// get the load averages
		if err := app.getLoadAvgs(ctx); err != nil {
			log.Error().Err(err).Msg("Could not get load averages.")
		}
		// send our data
		if err := hass.PublishState(app, client); err != nil {
			log.Error().Err(err).Msg("Failed to publish state.")
		}
	}

	// we use a helper function, PollSensors, that takes care of setting up a
	// ticker to run every minute (with a little bit of jitter) to send our data
	helpers.PollSensors(ctx, sendState, time.Minute, time.Second*5)
}

// Clear is the function that is called by the agent when we run the agent's
// clear command. It will remove all the sensors, button and their data from
// Home Assistant.
func Clear(_ context.Context, client hass.MQTTClient) {
	log.Info().Msgf("Clearing %s app data from Home Assistant.", appName)
	app := newExampleApp()

	if err := hass.Unpublish(app, client); err != nil {
		log.Error().Err(err).Msg("Failed to clear app data from Home Assistant.")
	}
	if err := app.config.UnRegister(appName); err != nil {
		log.Error().Err(err).Msg("Failed to remove registration.")
		return
	}
}

// buttonCallback is our callback function that is run when somebody presses the
// button we created in Home Assistant. In this case, we use xdg-open to open
// the Home Assistant homepage. But we could do other things, or many things at
// once and/or react based on the response data we got (the MQTT.Message
// parameter).
func buttonCallback(_ MQTT.Client, _ MQTT.Message) {
	if err := exec.Command("xdg-open", "https://home-assistant.io").Run(); err != nil {
		log.Error().Err(err).Msg("Could not execute xdg-open.")
	}
}
