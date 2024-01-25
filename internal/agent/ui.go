// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import ui "github.com/joshuar/go-hass-anything/v2/internal/agent/ui/bubbletea"

type AgentUI interface {
	Run()
	ShowConfiguration(ui.Agent)
}
