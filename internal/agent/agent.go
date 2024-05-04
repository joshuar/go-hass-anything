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

	ui "github.com/joshuar/go-hass-anything/v9/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/v9/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v9/pkg/preferences"
)

//go:generate go run ../../tools/appgenerator/run.go arg1
var (
	AppList []App
)

type agent struct {
	ui      AgentUI
	done    chan struct{}
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

func NewAgent(id, name string) *agent {
	a := &agent{
		id:   id,
		name: name,
	}
	a.ui = ui.NewBubbleTeaUI()
	return a
}

func (a *agent) AppName() string {
	return a.name
}

func (a *agent) AppID() string {
	return a.id
}

func (a *agent) AppVersion() string {
	return a.version
}

func (a *agent) Stop() {
	close(a.done)
}

func (a *agent) Configure() {
	a.ui.ShowConfiguration()
	a.ui.Run()
}

func Run(ctx context.Context) {
	var subscriptions []*mqtt.Subscription
	var configs []*mqtt.Msg
	for _, app := range AppList {
		configs = append(configs, app.Configuration()...)
		subscriptions = append(subscriptions, app.Subscriptions()...)
	}

	prefs, err := preferences.LoadPreferences()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load preferences.")
	}

	client, err := mqtt.NewClient(ctx, prefs, subscriptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to broker.")
	}

	runApps(ctx, client, AppList)
}

func ClearApps(ctx context.Context) {
	prefs, err := preferences.LoadPreferences()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load preferences.")
	}

	client, err := mqtt.NewClient(ctx, prefs, nil)
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
			publishAppConfigs(a, client)
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

func publishAppConfigs(app App, client *mqtt.Client) {
	log.Debug().Str("app", app.Name()).Msg("Publishing configs.")
	if err := client.Publish(app.Configuration()...); err != nil {
		log.Error().Err(err).Str("app", app.Name()).Msg("Failed to publish configuration messages.")
	}
}
