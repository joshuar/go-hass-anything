// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

var (
	// PreferencesDir is the default path under which the preferences are
	// written. It can be overridden in the Save/Load functions as needed.
	PreferencesDir = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")
	// PreferencesFile is the default filename used for storing the preferences
	// on disk. While it can be overridden, this is usually unnecessary.
	PreferencesFile    = "mqtt-config.toml"
	defaultPreferences = Preferences{
		Server:   "tcp://localhost:1883",
		User:     "",
		Password: "",
	}
	AppRegistryDir = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything", "appregistry")
)

// Preferences is a struct containing the general preferences for either the
// agent or for any code that imports go-hass-anything as a package. Currently,
// these preferences are for MQTT connectivity selection.
type Preferences struct {
	// Server is the URL for the MQTT server.
	Server string `toml:"mqttserver"`
	// User is the username for connecting to the server (optional).
	User string `toml:"mqttuser,omitempty"`
	// Password is the password for connecting to the server (optional).
	Password string `toml:"mqttpassword,omitempty"`
}

func (p *Preferences) MQTTServer() string {
	return p.Server
}

func (p *Preferences) MQTTUser() string {
	return p.User
}

func (p *Preferences) MQTTPassword() string {
	return p.Password
}

// Pref is a functional type for applying a value to a particular preference.
type Pref func(*Preferences)

// MQTTServer is the functional preference that sets the MQTTServer preference
// to the specified value.
func MQTTServer(server string) Pref {
	return func(args *Preferences) {
		args.Server = server
	}
}

// MQTTUser is the functional preference that sets the MQTTUser preference
// to the specified value.
func MQTTUser(user string) Pref {
	return func(args *Preferences) {
		args.User = user
	}
}

// MQTTPassword is the functional preference that sets the MQTTPassword preference
// to the specified value.
func MQTTPassword(password string) Pref {
	return func(args *Preferences) {
		args.Password = password
	}
}

// SavePreferences writes the given preferences to disk under the specified
// path. If the path is "", the preferences are saved to the file specified
// by PreferencesFile under the location specified by ConfigBasePath.
func SavePreferences(path string, setters ...Pref) error {
	if path == "" {
		path = PreferencesDir
	}
	file := filepath.Join(path, PreferencesFile)
	if err := checkPath(path); err != nil {
		return err
	}

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
		path = PreferencesDir
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

func Register(path, app string) error {
	if path == "" {
		path = AppRegistryDir
	}
	file := filepath.Join(path, app)
	if err := checkPath(path); err != nil {
		return err
	}

	if fs, err := os.Create(file); err != nil {
		return err
	} else {
		return fs.Close()
	}
}

func UnRegister(path, app string) error {
	if path == "" {
		path = AppRegistryDir
	}
	file := filepath.Join(path, app)
	return os.Remove(file)
}

func IsRegistered(path, app string) bool {
	if path == "" {
		path = AppRegistryDir
	}
	file := filepath.Join(path, app)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func checkPath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
