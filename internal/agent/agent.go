// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/joshuar/go-hass-anything/v9/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v9/pkg/preferences"
)

//go:generate go run ../../tools/appgenerator/run.go arg1
var (
	AppList []App
)

type Agent struct {
	done    chan struct{}
	prefs   *preferences.Preferences
	id      string
	name    string
	version string
}

// App represents an app that the agent can run. All apps have the following
// methods, which define how the app should be configured, current states of its
// entities and any subscriptions it wants to watch.
type App interface {
	// Name() is an identifier for the app, used for logging in the agent.
	Name() string
	// Configuration() returns the messages needed to tell Home Assistant how to
	// configure the app and its entities.
	Configuration() []*mqtt.Msg
	// States() are the messages that reflect the app's current state of the
	// entities of the app.
	States() []*mqtt.Msg
	// Subscriptions() are the topics on which the app wants to subscribe and
	// execute a callback in response to a message on that topic.
	Subscriptions() []*mqtt.Subscription
	// Update() is a function that is run at least once by the agent and will
	// usually contain the logic to update the states of all the apps entities.
	// It may be run multiple times, if the app is also considered a polling
	// app. See the definition for PollingApp for details.
	Update(ctx context.Context) error
}

// PollingApp represents an app that should be polled for updates on some
// interval. When an app satisfies this interface, the agent will configure a
// goroutine to run the apps Update() function and publish its States().
type PollingApp interface {
	// PollConfig defines the interval on which the app should be polled and its
	// states updated. A jitter should be defined, that is much less than the
	// interval, to add a small variation to the interval to avoid any
	// "thundering herd" problems.
	PollConfig() (interval, jitter time.Duration)
}

// EventsApp represents an app that will update its States() in response to some
// event(s) it is monitoring. When an app satisfies this interface, the agent
// will configure a goroutine to watch a channel of messages the app sends when
// an event occurs, which will be published to MQTT.
type EventsApp interface {
	// MsgCh is a channel of messages that the app generates when some internal
	// event occurs and a new message should be published to MQTT.
	MsgCh() chan *mqtt.Msg
}

//nolint:exhaustruct
func NewAgent(id, name string) *Agent {
	agent := &Agent{
		id:   id,
		name: name,
	}

	prefs, err := preferences.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Could not fetch agent preferences.")
	}

	agent.prefs = prefs

	return agent
}

func (a *Agent) AppName() string {
	return a.name
}

func (a *Agent) AppID() string {
	return a.id
}

func (a *Agent) AppVersion() string {
	return a.version
}

func (a *Agent) Stop() {
	close(a.done)
}

func (a *Agent) Name() string {
	return a.name
}

// GetPreferences returns the agent preferences.
func (a *Agent) GetPreferences() *preferences.Preferences {
	return a.prefs
}

// SetPreferences sets the agent preferences to the given preferences. If the
// preferences cannot be saved, a non-nil error is returned.
func (a *Agent) SetPreferences(prefs *preferences.Preferences) error {
	a.prefs = prefs

	return a.prefs.Save()
}

// Configure will start a terminal UI to adjust agent preferences and likewise for any
// apps that have user-configurable preferences.
func (a *Agent) Configure() {
	// Show a terminal UI to configure the agent preferences.
	if err := ShowPreferences(a); err != nil {
		log.Warn().Err(err).Msg("Problem occurred configuring agent.")
	}
	// For any apps that satisfy the Preferences interface, meaning they have
	// configurable preferences, show a terminal UI to configure them.
	for _, app := range AppList {
		if prefs, ok := app.(Preferences); ok {
			if err := ShowPreferences(prefs); err != nil {
				log.Warn().Err(err).Str("app", prefs.Name()).Msg("Problem occurred configuring app.")
			}
		}
	}
}

func Run(ctx context.Context) {
	var subscriptions []*mqtt.Subscription

	var configs []*mqtt.Msg

	for _, app := range AppList {
		configs = append(configs, app.Configuration()...)
		subscriptions = append(subscriptions, app.Subscriptions()...)
	}

	prefs, err := preferences.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load preferences.")
	}

	client, err := mqtt.NewClient(ctx, prefs, subscriptions, configs)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to broker.")
	}

	runApps(ctx, client, AppList)
}

func ClearApps(ctx context.Context) {
	prefs, err := preferences.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load preferences.")
	}

	client, err := mqtt.NewClient(ctx, prefs, nil, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to broker.")
	}

	for _, app := range AppList {
		log.Debug().Str("app", app.Name()).Msg("Removing configuration for app.")
		if err := client.Unpublish(app.Configuration()...); err != nil {
			log.Error().Err(err).Str("app", app.Name()).Msg("Could not remove configuration for app.")

			continue
		}
	}
}

func runApps(ctx context.Context, client *mqtt.Client, apps []App) {
	var wg sync.WaitGroup

	for _, app := range apps {
		log.Debug().Str("app", app.Name()).Msg("Running app.")
		wg.Add(1)
		go func(a App) {
			defer wg.Done()

			if app, ok := a.(PollingApp); ok {
				interval, jitter := app.PollConfig()
				log.Info().Dur("interval", interval).Str("app", a.Name()).Msg("Running loop to poll app updates.")
				wg.Add(1)

				go func() {
					defer wg.Wait()
					poll(
						ctx,
						func() {
							updateApp(ctx, a)
							publishAppStates(a, client)
						},
						interval,
						jitter,
					)
				}()
			}

			if app, ok := a.(EventsApp); ok {
				updateApp(ctx, a)
				wg.Add(1)

				go func() {
					defer wg.Done()
					log.Info().Str("app", a.Name()).Msg("Listening for message events from app to publish.")

					for {
						select {
						case msg := <-app.MsgCh():
							if err := client.Publish(msg); err != nil {
								log.Error().Err(err).Str("app", a.Name()).Msg("Failed to publish state messages.")
							}
						case <-ctx.Done():
							return
						}
					}
				}()
			}
		}(app)
	}

	wg.Wait()
}

func updateApp(ctx context.Context, app App) {
	log.Debug().Str("app", app.Name()).Msg("Updating.")

	if err := app.Update(ctx); err != nil {
		log.Warn().Err(err).Str("app", app.Name()).Msg("App failed to update.")
	}
}

func publishAppStates(app App, client *mqtt.Client) {
	log.Debug().Str("app", app.Name()).Msg("Publishing states.")

	if err := client.Publish(app.States()...); err != nil {
		log.Error().Err(err).Str("app", app.Name()).Msg("Failed to publish state messages.")
	}
}
