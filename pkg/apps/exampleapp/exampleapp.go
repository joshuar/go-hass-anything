// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package exampleapp

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/carlmjohnson/requests"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/load"

	"github.com/joshuar/go-hass-anything/v5/pkg/apps/helpers"
	"github.com/joshuar/go-hass-anything/v5/pkg/hass"
	"github.com/joshuar/go-hass-anything/v5/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v5/pkg/preferences"
	"github.com/joshuar/go-hass-anything/v5/pkg/web"
)

const (
	appName        = "exampleApp"
	weatherURL     = "http://wttr.in/?format=j1"
	weatherURLpref = "weatherURL"
)

var defaultPrefs = map[string]any{
	"weatherURL": weatherURL,
}

type exampleApp struct {
	loadData    *load.AvgStat
	config      *preferences.AppPreferences
	weatherData []byte
}

// New sets up our example app. We make use of the preference loading/saving in
// the agent to provide a file for our app preferences at
// ~/.config/go-hass-anything/exampleApp-preferences.toml. We can store whatever
// preferences our app needs in this file by providing a map[string]any that
// maps preferences to values.
func New(_ context.Context) (*exampleApp, error) {
	app := &exampleApp{}
	// load our app config. if we don't have a config, set some defaults
	p, err := preferences.LoadAppPreferences(app.Name())
	if os.IsNotExist(err) {
		log.Info().Msgf("Setting default weather service to %s", weatherURL)
		p.Prefs = defaultPrefs
		err := preferences.SaveAppPreferences(app.Name(), preferences.SetAppPreferences(p.Prefs))
		if err != nil {
			log.Warn().Err(err).Msg("Could not save app preferences.")
		}
	}
	app.config = p
	return app, nil
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
	r := <-web.ExecuteRequest(ctx, a)
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
	var ok bool
	if serviceURL, ok = a.config.Prefs[weatherURLpref].(string); !ok {
		log.Error().Msg("Could not retrieve weather service URL from config.")
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
func (a *exampleApp) Name() string {
	return appName
}

// Configuration is called when our app is first registered in Home Assistant and
// will return configuration messages for the data our app will send/receive.
func (a *exampleApp) Configuration() []*mqtt.Msg {
	var msgs []*mqtt.Msg
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
			WithDeviceInfo(deviceInfo).
			WithOriginInfo(originInfo).
			WithStateClassMeasurement().
			WithDeviceClass("temperature").
			WithUnits("°C").
			WithValueTemplate("{{ value_json.current_condition[0].temp_C }}"))

	// we have three sensors for the loadavgs
	for _, l := range []string{"1", "5", "15"} {
		entities = append(entities,
			hass.NewEntityByID("example_app_load"+l, appName).
				AsSensor().
				WithDeviceInfo(deviceInfo).
				WithOriginInfo(originInfo).
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

// States is called when we want to send our sensor data to Home Assistant.
func (a *exampleApp) States() []*mqtt.Msg {
	var msgs []*mqtt.Msg

	// we retrieve the weather data and send that as the weather sensor state
	msgs = append(msgs,
		mqtt.NewMsg(mqtt.DiscoveryPrefix+"/sensor/"+appName+"/example_app_weather_temp/state", a.weatherData))

	// we retrieve our load avgs
	for _, loads := range []string{"1", "5", "15"} {
		id := "example_app_load" + loads
		var l float64
		switch loads {
		case "1":
			l = a.loadData.Load1
		case "5":
			l = a.loadData.Load5
		case "15":
			l = a.loadData.Load15
		}
		msgs = append(msgs,
			mqtt.NewMsg(
				mqtt.DiscoveryPrefix+"/sensor/"+appName+"/"+id+"/state",
				json.RawMessage(strconv.FormatFloat(l, 'f', -1, 64))))
	}

	return msgs
}

// Subscriptions is called once to register the callbacks we will use when Home
// Assistant sends back messages on any command topics.
func (a *exampleApp) Subscriptions() []*mqtt.Subscription {
	var msgs []*mqtt.Subscription

	// we add our callback for our button
	msgs = append(msgs, &mqtt.Subscription{
		Topic:    mqtt.DiscoveryPrefix + "/button/" + appName + "/example_app_button/toggle",
		Callback: buttonCallback,
	})
	return msgs
}

// Run is the function that the agent calls to start our app. In it, we create
// our app struct, register our app (if needed), listen for our button press,
// then set up a loop to send our sensor data.
func (a *exampleApp) Run(ctx context.Context, client hass.MQTTClient) error {
	log.Info().Str("appName", appName).Msg("Starting app.")

	// add our button subscription
	if err := hass.Subscribe(a, client); err != nil {
		log.Error().Err(err).Msg("Could not activate subscriptions.")
	}

	// create a function we will use whenever we want to send our app state
	sendState := func() {
		// get the weather
		if err := a.getWeather(ctx); err != nil {
			log.Error().Err(err).Msg("Could not get weather data.")
		}
		// get the load averages
		if err := a.getLoadAvgs(ctx); err != nil {
			log.Error().Err(err).Msg("Could not get load averages.")
		}
		// send our data
		if err := hass.PublishState(a, client); err != nil {
			log.Error().Err(err).Msg("Failed to publish state.")
		}
	}

	// we use a helper function, PollSensors, that takes care of setting up a
	// ticker to run every minute (with a little bit of jitter) to send our data
	helpers.PollSensors(ctx, sendState, time.Minute, time.Second*5)
	return nil
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
