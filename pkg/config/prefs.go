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
	PreferencesFile = "config.toml"
)

type Preferences struct {
	MQTTServer   string `toml:"mqttserver"`
	MQTTUser     string `toml:"mqttuser,omitempty"`
	MQTTPassword string `toml:"mqttpassword,omitempty"`
}

type Pref func(*Preferences)

func MQTTServer(server string) Pref {
	return func(args *Preferences) {
		args.MQTTServer = server
	}
}

func MQTTUser(user string) Pref {
	return func(args *Preferences) {
		args.MQTTUser = user
	}
}

func MQTTPassword(password string) Pref {
	return func(args *Preferences) {
		args.MQTTPassword = password
	}
}

func SavePreferences(path string, setters ...Pref) error {
	if path == "" {
		path = ConfigBasePath
	}
	file := filepath.Join(path, PreferencesFile)
	checkPath(path)

	args := &Preferences{
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
	err = os.WriteFile(file, b, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func LoadPreferences(path string) (*Preferences, error) {
	if path == "" {
		path = ConfigBasePath
	}
	file := filepath.Join(path, PreferencesFile)

	p := &Preferences{}
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func checkPath(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Debug().Err(err).Msgf("Failed to create new config directory %s.", path)
		} else {
			log.Debug().Msgf("Created new config directory %s.", path)
		}
	}
}
