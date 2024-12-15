// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/joshuar/go-hass-anything/v12/internal/logging"
	"github.com/joshuar/go-hass-anything/v12/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v12/pkg/preferences"
)

//go:generate go run ../../tools/appgenerator/main.go
var (
	AppList []App
)

// Agent is a structure holding the internal agent data.
type Agent struct {
	done    chan struct{}
	logger  *slog.Logger
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

// AppWithPreferences represents an app that has preferences that can be
// configured by the user.
type AppWithPreferences interface {
	App
	// Preferences returns the AppPreferences map of preferences for the app.
	// This is passed to the UI code to facilitate generating a form to enter
	// the preferences when the agent runs its configure command. If the
	// preferences cannot be returned, a non-nil error will be returned.
	DefaultPreferences() preferences.AppPreferences
}

// PollingApp represents an app that should be polled for updates on some
// interval. When an app satisfies this interface, the agent will configure a
// goroutine to run the apps Update() function and publish its States().
type PollingApp interface {
	App
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
	App
	// MsgCh is a channel of messages that the app generates when some internal
	// event occurs and a new message should be published to MQTT.
	MsgCh() chan *mqtt.Msg
}

//nolint:exhaustruct
func NewAgent(ctx context.Context, id, name string) *Agent {
	agent := &Agent{
		id:     id,
		name:   name,
		logger: logging.FromContext(ctx).With(slog.String("source", id)),
	}

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

// Configure will start a terminal UI to adjust agent preferences and likewise for any
// apps that have user-configurable preferences.
func (a *Agent) Configure() {
	// Get the agent preferences.
	if err := preferences.Load(); err != nil {
		a.logger.Warn("Could not load agent preferences.",
			slog.Any("error", err))
	}
	// Show a terminal UI to configure the agent preferences.
	if err := ShowPreferences(a.AppName(), preferences.Agent); err != nil {
		a.logger.Warn("Could not display agent preferences.",
			slog.Any("error", err))
	}
	// Save agent preferences.
	if err := preferences.Save(); err != nil {
		a.logger.Warn("Could not save agent preferences.",
			slog.Any("error", err))
	}
	// For any apps that satisfy the Preferences interface, meaning they have
	// configurable preferences, show a terminal UI to configure them.
	for _, app := range AppList {
		a.logger.Debug("Checking and configuring app preferences",
			slog.String("app", app.Name()))

		app, ok := app.(AppWithPreferences)
		if !ok {
			continue
		}

		appPrefs, err := preferences.LoadApp(app)
		if err != nil {
			a.logger.Warn("Could not configure app.",
				slog.String("app", app.Name()),
				slog.Any("error", err))

			continue
		}

		if err := ShowPreferences(app.Name(), appPrefs); err != nil {
			a.logger.Warn("Could not configure app.",
				slog.String("app", app.Name()),
				slog.Any("error", err))
		}

		if err := preferences.SaveApp(app.Name(), appPrefs); err != nil {
			a.logger.Warn("Could not configure app.",
				slog.String("app", app.Name()),
				slog.Any("error", err))
		}
	}
}

func Run(ctx context.Context) error {
	var (
		subscriptions []*mqtt.Subscription
		configs       []*mqtt.Msg
	)

	if err := preferences.Load(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	for _, app := range AppList {
		configs = append(configs, app.Configuration()...)
		subscriptions = append(subscriptions, app.Subscriptions()...)
	}

	client, err := mqtt.NewClient(ctx, preferences.Agent, subscriptions, configs)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	runApps(ctx, client, AppList)

	return nil
}

func ClearApps(ctx context.Context) error {
	if err := preferences.Load(); err != nil {
		return fmt.Errorf("clear: %w", err)
	}

	client, err := mqtt.NewClient(ctx, preferences.Agent, nil, nil)
	if err != nil {
		return fmt.Errorf("clear: %w", err)
	}

	for _, app := range AppList {
		logging.FromContext(ctx).Debug("Removing configuration from MQTT for app.",
			slog.String("app", app.Name()))

		if err := client.Unpublish(ctx, app.Configuration()...); err != nil {
			logging.FromContext(ctx).Warn("Could not remove configuration from MQTT for app.",
				slog.String("app", app.Name()),
				slog.Any("error", err))

			continue
		}
	}

	return nil
}

func runApps(ctx context.Context, client *mqtt.Client, apps []App) {
	var wg sync.WaitGroup

	logger := logging.FromContext(ctx)

	for _, app := range apps {
		logger.Debug("Running app.",
			slog.String("app", app.Name()))

		switch app := app.(type) {
		case PollingApp:
			wg.Add(1)

			go func() {
				defer wg.Done()
				runPollingApp(ctx, client, logger, app)
			}()
		case EventsApp:
			wg.Add(1)

			go func() {
				defer wg.Done()
				runEventsApp(ctx, client, logger, app)
			}()
		default:
			updateApp(ctx, app)
			publishAppStates(ctx, app, client)
		}
	}

	wg.Wait()
}

func updateApp(ctx context.Context, app App) {
	logger := logging.FromContext(ctx)

	logger.Debug("Updating app.",
		slog.String("app", app.Name()))

	if err := app.Update(ctx); err != nil {
		logger.Warn("Failed to update app.",
			slog.String("app", app.Name()),
			slog.Any("error", err))
	}
}

func publishAppStates(ctx context.Context, app App, client *mqtt.Client) {
	logger := logging.FromContext(ctx)

	logger.Debug("Publishing app states.",
		slog.String("app", app.Name()))

	if err := client.Publish(ctx, app.States()...); err != nil {
		logger.Warn("Failed to publish app states.",
			slog.String("app", app.Name()),
			slog.Any("error", err))
	}
}

func runPollingApp(ctx context.Context, client *mqtt.Client, logger *slog.Logger, app PollingApp) {
	interval, jitter := app.PollConfig()

	logger.Info("Running loop to poll app for updates.",
		slog.String("app", app.Name()))

	err := poll(
		ctx,
		func() {
			updateApp(ctx, app)
			publishAppStates(ctx, app, client)
		},
		interval,
		jitter,
	)
	if err != nil {
		logger.Error("Failed to poll app for updates.",
			slog.String("app", app.Name()),
			slog.Any("error", err))
	}
}

func runEventsApp(ctx context.Context, client *mqtt.Client, logger *slog.Logger, app EventsApp) {
	updateApp(ctx, app)

	logger.Info("Listening for message events from app.",
		slog.String("app", app.Name()))

	for {
		select {
		case msg := <-app.MsgCh():
			if err := client.Publish(ctx, msg); err != nil {
				logger.Error("Failed to publish state messages for app.",
					slog.String("app", app.Name()),
					slog.Any("error", err))
			}
		case <-ctx.Done():
			return
		}
	}
}
