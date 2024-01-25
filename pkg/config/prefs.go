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
	// ConfigBasePath is the default path under which the preferences are
	// written. It can be overridden in the Save/Load functions as needed.
	ConfigBasePath = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")
	// PreferencesFile is the default filename used for storing the preferences
	// on disk. While it can be overridden, this is usually unnecessary.
	PreferencesFile    = "config.toml"
	defaultPreferences = Preferences{
		MQTTServer:   "localhost:1883",
		MQTTUser:     "",
		MQTTPassword: "",
	}
)

// Preferences is a struct containing the general preferences for either the
// agent or for any code that imports go-hass-anything as a package. Currently,
// these preferences are for MQTT connectivity selection.
type Preferences struct {
	// MQTTServer is the URL for the MQTT server.
	MQTTServer string `toml:"mqttserver"`
	// MQTTUser is the username for connecting to the server (optional).
	MQTTUser string `toml:"mqttuser,omitempty"`
	// MQTTPassword is the password for connecting to the server (optional).
	MQTTPassword string `toml:"mqttpassword,omitempty"`
}

// Pref is a functional type for applying a value to a particular preference.
type Pref func(*Preferences)

// MQTTServer is the functional preference that sets the MQTTServer preference
// to the specified value.
func MQTTServer(server string) Pref {
	return func(args *Preferences) {
		args.MQTTServer = server
	}
}

// MQTTUser is the functional preference that sets the MQTTUser preference
// to the specified value.
func MQTTUser(user string) Pref {
	return func(args *Preferences) {
		args.MQTTUser = user
	}
}

// MQTTPassword is the functional preference that sets the MQTTPassword preference
// to the specified value.
func MQTTPassword(password string) Pref {
	return func(args *Preferences) {
		args.MQTTPassword = password
	}
}

// SavePreferences writes the given preferences to disk under the specified
// path. If the path is "", the preferences are saved to the file specified
// by PreferencesFile under the location specified by ConfigBasePath.
func SavePreferences(path string, setters ...Pref) error {
	if path == "" {
		path = ConfigBasePath
	}
	file := filepath.Join(path, PreferencesFile)
	checkPath(path)

	args := defaultPreferences
	for _, setter := range setters {
		setter(&args)
	}

	b, err := toml.Marshal(&args)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, b, 0o600)
	if err != nil {
		return err
	}
	return nil
}

// LoadPreferences retrives all Preferences from disk at the given path. If the
// path is "", the preferences are loaded from the file specified by
// PreferencesFile under the location specified by ConfigBasePath.
func LoadPreferences(path string) (*Preferences, error) {
	if path == "" {
		path = ConfigBasePath
	}
	file := filepath.Join(path, PreferencesFile)

	p := defaultPreferences
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
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
