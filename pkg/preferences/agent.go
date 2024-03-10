// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package preferences

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
)

var (
	// preferencesDir is the default path under which the preferences are
	// written. It can be overridden in the Save/Load functions as needed.
	preferencesDir = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")
	// preferencesFile is the default filename used for storing the preferences
	// on disk. While it can be overridden, this is usually unnecessary.
	preferencesFile = "mqtt-config.toml"
	// defaultServer is the default MQTT broker URL.
	defaultServer = "tcp://localhost:1883"
	// defaultTopicPrefix is the default prefix that is appended to topics.
	defaultTopicPrefix = "homeassistant"
)

// Preferences is a struct containing the general preferences for either the
// agent or for any code that imports go-hass-anything as a package. Currently,
// these preferences are for MQTT connectivity selection.
type Preferences struct {
	Server      string `toml:"mqttserver"`
	User        string `toml:"mqttuser,omitempty"`
	Password    string `toml:"mqttpassword,omitempty"`
	TopicPrefix string `toml:"topicprefix"`
}

// MQTTServer returns the current server set in the preferences.
func (p *Preferences) GetMQTTServer() string {
	return p.Server
}

// MQTTUser returns any username set in the preferences.
func (p *Preferences) GetMQTTUser() string {
	return p.User
}

// MQTTPassword returns any password set in the preferences.
func (p *Preferences) GetMQTTPassword() string {
	return p.Password
}

// TopicPrefix returns the topic prefix set in the preferences.
func (p *Preferences) GetTopicPrefix() string {
	return p.TopicPrefix
}

// Pref is a functional type for applying a value to a particular preference.
type Pref func(*Preferences)

// SetMQTTServer is the functional preference that sets the MQTT server preference
// to the specified value.
func SetMQTTServer(server string) Pref {
	return func(args *Preferences) {
		args.Server = server
	}
}

// SetMQTTUser is the functional preference that sets the MQTT user preference
// to the specified value.
func SetMQTTUser(user string) Pref {
	return func(args *Preferences) {
		args.User = user
	}
}

// SetMQTTPassword is the functional preference that sets the MQTT password preference
// to the specified value.
func SetMQTTPassword(password string) Pref {
	return func(args *Preferences) {
		args.Password = password
	}
}

// SetTopicPrefix is the functional preference that sets the topic prefix preference
// to the specified value.
func SetTopicPrefix(prefix string) Pref {
	return func(args *Preferences) {
		args.TopicPrefix = prefix
	}
}

// SavePreferences writes the given preferences to disk under the specified
// path. If the path is "", the preferences are saved to the file specified
// by PreferencesFile under the location specified by ConfigBasePath.
func SavePreferences(setters ...Pref) error {
	file := filepath.Join(preferencesDir, preferencesFile)
	if err := checkPath(preferencesDir); err != nil {
		return err
	}

	prefs, err := LoadPreferences()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, setter := range setters {
		setter(prefs)
	}

	return write(prefs, file)
}

// LoadPreferences retrives all Preferences from disk at the given path. If the
// path is "", the preferences are loaded from the file specified by
// PreferencesFile under the location specified by ConfigBasePath.
func LoadPreferences() (*Preferences, error) {
	file := filepath.Join(preferencesDir, preferencesFile)

	p := defaultPreferences()
	b, err := os.ReadFile(file)
	if err != nil {
		return p, err
	}
	err = toml.Unmarshal(b, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// SetPath will set the path to the preferences file to the given path. This
// function is optional and a default path is used when calling Load/Save if it
// is not called.
func SetPath(path string) {
	preferencesDir = path
}

// SetFile will set the filename of the preferences file to the given name. This
// function is optional and a default filename is used when calling Load/Save if
// it is not called.
func SetFile(file string) {
	preferencesFile = file
}

func defaultPreferences() *Preferences {
	return &Preferences{
		Server:      defaultServer,
		User:        "",
		Password:    "",
		TopicPrefix: defaultTopicPrefix,
	}
}

func write(prefs any, file string) error {
	b, err := toml.Marshal(prefs)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, b, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func checkPath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Debug().Str("preferencesPath", path).
			Msg("Creating new preferences path.")
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
