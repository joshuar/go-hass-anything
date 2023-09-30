// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"

	ui "github.com/joshuar/go-hass-anything/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/internal/apps/exampleApp"
	"github.com/joshuar/go-hass-anything/pkg/config"
	viper "github.com/joshuar/go-hass-anything/pkg/config/viper"
	"github.com/joshuar/go-hass-anything/pkg/hass"
	"github.com/joshuar/go-hass-anything/pkg/mqtt"
	"github.com/rs/zerolog/log"
)

var (
	RunList   []func(context.Context, hass.MQTTClient)
	ClearList []func(context.Context, hass.MQTTClient)
)

func init() {
	RunList = append(RunList, exampleApp.Run)
	ClearList = append(ClearList, exampleApp.Clear)
}

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

func (a *agent) RunApps(ctx context.Context, mqtt *mqtt.MQTTClient) {
	appCh := make(chan func(context.Context, hass.MQTTClient), len(RunList))
	var wg sync.WaitGroup

	for i := 0; i < len(RunList); i++ {
		appCh <- RunList[i]
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for app := range appCh {
			app(ctx, mqtt)
		}
	}()

	close(appCh)
	wg.Wait()
}

func (a *agent) ClearApps(ctx context.Context, mqtt *mqtt.MQTTClient) {
	appCh := make(chan func(context.Context, hass.MQTTClient), len(ClearList))
	var wg sync.WaitGroup

	for i := 0; i < len(ClearList); i++ {
		appCh <- ClearList[i]
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for app := range appCh {
			app(ctx, mqtt)
		}
	}()

	close(appCh)
	wg.Wait()
}
