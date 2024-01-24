// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
)

var (
	ConfigBasePath  = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")
	PreferencesFile = filepath.Join(ConfigBasePath, "config.toml")
)

type AppPreferences struct {
	MQTTServer   string `toml:"mqttserver"`
	MQTTUser     string `toml:"mqttuser,omitempty"`
	MQTTPassword string `toml:"mqttpassword,omitempty"`
}

type Pref func(*AppPreferences)

func MQTTServer(server string) Pref {
	return func(args *AppPreferences) {
		args.MQTTServer = server
	}
}

func MQTTUser(user string) Pref {
	return func(args *AppPreferences) {
		args.MQTTUser = user
	}
}

func MQTTPassword(password string) Pref {
	return func(args *AppPreferences) {
		args.MQTTPassword = password
	}
}

func SavePreferences(setters ...Pref) error {
	args := &AppPreferences{
		MQTTServer:   "localhost:1883",
		MQTTUser:     "",
		MQTTPassword: "",
	}
	for _, setter := range setters {
		setter(args)
	}

	b, err := toml.Marshal(args)
	if err != nil {
		return err
	}
	err = os.WriteFile(PreferencesFile, b, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func LoadPreferences() (*AppPreferences, error) {
	p := &AppPreferences{}
	b, err := os.ReadFile(PreferencesFile)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

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
