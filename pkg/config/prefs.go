// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

var ConfigBasePath = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")

const (
	PrefMQTTServer   = "mqttserver"
	PrefMQTTUser     = "mqttuser"
	PrefMQTTPassword = "mqttpassword"
)

func init() {
	_, err := os.Stat(ConfigBasePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(ConfigBasePath, os.ModePerm); err != nil {
			log.Debug().Err(err).Msgf("Failed to create new config directory %s.", ConfigBasePath)
		} else {
			log.Debug().Msgf("Created new config directory %s.", ConfigBasePath)
		}
	}
}
