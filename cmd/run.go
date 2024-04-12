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

	"github.com/spf13/cobra"

	"github.com/joshuar/go-hass-anything/v7/internal/agent"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run go-hass-anything",
	Long:  `Run go-hass-anything and start sending all configured apps and their data to Home Assistant.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancelFunc := context.WithCancel(context.TODO())
		defer cancelFunc()

		agent.Run(ctx)
	},
}
