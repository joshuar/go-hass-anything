// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	ui "github.com/joshuar/go-hass-anything/v8/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/v8/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v8/pkg/preferences"
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

type App interface {
	Name() string
	Configuration() []*mqtt.Msg
	States() []*mqtt.Msg
	Subscriptions() []*mqtt.Subscription
	Run(ctx context.Context, client *mqtt.Client) error
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
	var appsToRun []mqtt.Device
	for _, app := range AppList {
		appsToRun = append(appsToRun, app)
	}

	prefs, err := preferences.LoadPreferences()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load preferences.")
	}

	client, err := mqtt.NewClient(ctx, prefs, appsToRun...)
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

	client, err := mqtt.NewClient(ctx, prefs)
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
			if err := a.Run(ctx, client); err != nil {
				log.Debug().Err(err).Msg("Error running app.")
			}
		}(app)
	}
	wg.Wait()
}
