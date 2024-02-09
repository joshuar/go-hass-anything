// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	ui "github.com/joshuar/go-hass-anything/v5/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/v5/pkg/hass"
	"github.com/joshuar/go-hass-anything/v5/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v5/pkg/preferences"
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
	Run(ctx context.Context, client hass.MQTTClient) error
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

func Run(ctx context.Context, client hass.MQTTClient) {
	prefs, err := preferences.LoadPreferences()
	if err != nil {
		log.Error().Err(err).Msg("Could not load preferences.")
		return
	}

	var appsToRun []App
	var registeredApps []string
	for _, app := range AppList {
		if prefs.IsRegistered(app.Name()) {
			appsToRun = append(appsToRun, app)
		} else {
			if err := hass.Register(app, client); err != nil {
				log.Error().Err(err).Str("app", app.Name()).Msg("Could not register app.")
				continue
			}
			appsToRun = append(appsToRun, app)
			registeredApps = append(registeredApps, app.Name())
			log.Info().Str("app", app.Name()).Msg("App registered.")
		}
	}
	if err := preferences.SavePreferences(preferences.RegisterApps(registeredApps...)); err != nil {
		log.Warn().Err(err).Msg("Failed to save registered apps in preferences.")
	}
	runApps(ctx, client, appsToRun)
}

func ClearApps(ctx context.Context, client hass.MQTTClient) {
	prefs, err := preferences.LoadPreferences()
	if err != nil {
		log.Error().Err(err).Msg("Could not load preferences.")
		return
	}
	var unRegisteredApps []string
	for _, app := range AppList {
		if prefs.IsRegistered(app.Name()) {
			if err := hass.UnRegister(app, client); err != nil {
				log.Error().Err(err).Str("app", app.Name()).Msg("Could not unregister app.")
				continue
			}
			unRegisteredApps = append(unRegisteredApps, app.Name())
		}
	}
	if err := preferences.SavePreferences(preferences.UnRegisterApps(unRegisteredApps...)); err != nil {
		log.Warn().Err(err).Msg("Failed to update registered apps in preferences.")
	}
}

func runApps(ctx context.Context, client hass.MQTTClient, apps []App) {
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
	log.Debug().Msg("here")
}
