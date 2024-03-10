// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	ui "github.com/joshuar/go-hass-anything/v6/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/v6/pkg/hass"
	"github.com/joshuar/go-hass-anything/v6/pkg/mqtt"
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
	var appsToRun []hass.MQTTDevice
	for _, app := range AppList {
		appsToRun = append(appsToRun, app)
	}
	if err := hass.Register(client, appsToRun...); err != nil {
		log.Error().Err(err).Msg("Could not register.")
	}
	runApps(ctx, client, AppList)
}

func ClearApps(ctx context.Context, client hass.MQTTClient) {
	for _, app := range AppList {
		if err := hass.UnRegister(client, app); err != nil {
			log.Error().Err(err).Str("app", app.Name()).Msg("Could not unregister app.")
			continue
		}
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
