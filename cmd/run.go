// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	"github.com/joshuar/go-hass-anything/internal/agent"
	"github.com/joshuar/go-hass-anything/pkg/mqtt"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run go-hass-anything",
	Long:  `Run go-hass-anything and start sending all configured apps and their data to Home Assistant.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		agent := agent.NewAgent(AgentID, AgentName)
		mqtt, err := mqtt.NewMQTTClient()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not connect to broker.")
		}
		agent.RunApps(ctx, mqtt)
	},
}
