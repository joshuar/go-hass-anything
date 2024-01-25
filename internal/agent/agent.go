// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"sync"

	ui "github.com/joshuar/go-hass-anything/v2/internal/agent/ui/bubbletea"
	"github.com/joshuar/go-hass-anything/v2/pkg/hass"
	"github.com/joshuar/go-hass-anything/v2/pkg/mqtt"
)

//go:generate go run ../../tools/appgenerator/run.go arg1

var (
	RunList   []func(context.Context, hass.MQTTClient)
	ClearList []func(context.Context, hass.MQTTClient)
)

type agent struct {
	ui      AgentUI
	done    chan struct{}
	id      string
	name    string
	version string
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

func (a *agent) RunApps(ctx context.Context, client *mqtt.Client) {
	a.doApps(ctx, client, RunList)
}

func (a *agent) ClearApps(ctx context.Context, client *mqtt.Client) {
	a.doApps(ctx, client, ClearList)
}

func (a *agent) doApps(ctx context.Context, client *mqtt.Client, appList []func(context.Context, hass.MQTTClient)) {
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
