// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/joshuar/go-hass-anything/v3/internal/agent"
)

var configCmd = &cobra.Command{
	Use:   "configure",
	Short: "Display configuration",
	Long:  `Open the configuration UI to make changes`,
	Run: func(cmd *cobra.Command, args []string) {
		agent := agent.NewAgent(AgentID, AgentName)
		agent.Configure()
	},
}
