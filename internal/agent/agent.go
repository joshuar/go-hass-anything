// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"

	ui "github.com/joshuar/go-hass-anything/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/pkg/config"
	viper "github.com/joshuar/go-hass-anything/pkg/config/viper"
	"github.com/joshuar/go-hass-anything/pkg/hass"
	"github.com/joshuar/go-hass-anything/pkg/mqtt"
	"github.com/rs/zerolog/log"
)

//go:generate go run ../../tools/appgenerator/run.go arg1

var (
	RunList   []func(context.Context, hass.MQTTClient)
	ClearList []func(context.Context, hass.MQTTClient)
)

type agent struct {
	Config  config.Agent
	ui      AgentUI
	done    chan struct{}
	id      string
	name    string
	version string
}

func NewAgent(id, name string) *agent {
	var err error
	a := &agent{
		id:   id,
		name: name,
	}
	a.ui = ui.NewBubbleTeaUI()
	if a.Config, err = viper.Load(""); err != nil {
		log.Warn().Err(err).Msg("No agent config found.")
	}
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

func (a *agent) GetConfig(key string, value interface{}) error {
	return a.Config.Get(key, value)
}

func (a *agent) SetConfig(key string, value interface{}) error {
	return a.Config.Set(key, value)
}

func (a *agent) Configure() {
	a.ui.ShowConfiguration(a)
	a.ui.Run()
}

func (a *agent) RunApps(ctx context.Context, client *mqtt.MQTTClient) {
	a.doApps(ctx, client, RunList)
}

func (a *agent) ClearApps(ctx context.Context, client *mqtt.MQTTClient) {
	a.doApps(ctx, client, ClearList)
}

func (a *agent) doApps(ctx context.Context, client *mqtt.MQTTClient, appList []func(context.Context, hass.MQTTClient)) {
	var wg sync.WaitGroup
	for _, app := range appList {
		wg.Add(1)
		go func(a func(context.Context, hass.MQTTClient)) {
			defer wg.Done()
			a(ctx, client)
		}(app)
	}
	wg.Wait()
}
