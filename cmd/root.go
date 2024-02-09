// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/joshuar/go-hass-anything/v5/internal/logging"
)

const (
	AgentID   = "com.github.joshuar.go-hass-anything"
	AgentName = "Go HASS Anything"
)

var (
	traceFlag   bool
	debugFlag   bool
	profileFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-hass-anything",
	Short: "Send any data to Home Assistant with Go and MQTT.",
	Long:  `Send any data to Home Assistant with Go and MQTT. If you can access it with a Go program, go-hass-anything can send it to Home Assistant for you.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging.SetLogging(traceFlag, debugFlag, profileFlag)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not run command.")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&traceFlag, "trace", false,
		"trace output (default is false)")
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false,
		"debug output (default is false)")
	rootCmd.PersistentFlags().BoolVar(&profileFlag, "profile", false,
		"enable profiling (default is false)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(clearCmd)
}
