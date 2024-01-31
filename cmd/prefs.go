// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// prefsCmd represents the prefs command
var prefsCmd = &cobra.Command{
	Use:   "prefs",
	Short: "Get/Set agent preferences.",
	Long: `Get or Set agent preferences For example:

go-hass-anything get server # gets the MQTT server
go-hass-anything set server some.host:port # sets the MQTT server 	
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("prefs called")
	},
}

func init() {
	rootCmd.AddCommand(prefsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// prefsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// prefsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
