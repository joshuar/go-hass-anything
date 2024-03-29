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

	"github.com/joshuar/go-hass-anything/v6/internal/agent"
	"github.com/joshuar/go-hass-anything/v6/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v6/pkg/preferences"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear app state from Home Assistant.",
	Long:  `Publish messages to MQTT to instruct Home Assistant to remove apps and their data.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		prefs, err := preferences.LoadPreferences()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not load preferences.")
		}

		mqtt, err := mqtt.NewMQTTClient(ctx, prefs)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not connect to broker.")
		}
		agent.ClearApps(ctx, mqtt)
	},
}
