// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package preferences

import (
	"os"
	"path/filepath"
	"slices"

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
)

// Preferences is a struct containing the general preferences for either the
// agent or for any code that imports go-hass-anything as a package. Currently,
// these preferences are for MQTT connectivity selection.
type Preferences struct {
	Server         string   `toml:"mqttserver"`
	User           string   `toml:"mqttuser,omitempty"`
	Password       string   `toml:"mqttpassword,omitempty"`
	RegisteredApps []string `toml:"registeredapps,omitempty"`
}

// MQTTServer returns the current server set in the preferences.
func (p *Preferences) MQTTServer() string {
	return p.Server
}

// MQTTUser returns any username set in the preferences.
func (p *Preferences) MQTTUser() string {
	return p.User
}

// MQTTPassword returns any password set in the preferences.
func (p *Preferences) MQTTPassword() string {
	return p.Password
}

// IsRegistered will check whether the given app has been recorded as registered
// in the preferences.
func (p *Preferences) IsRegistered(app string) bool {
	return slices.Contains(p.RegisteredApps, app)
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

// RegisterApps is the functional preference that appends the list of given apps
// to the existing registered apps in the preferences.
func RegisterApps(apps ...string) Pref {
	return func(args *Preferences) {
		args.RegisteredApps = append(args.RegisteredApps, apps...)
	}
}

// UnRegisterApps is the functional preference that will remove the list of
// given apps from the registered apps in the preferences.
func UnRegisterApps(apps ...string) Pref {
	return func(args *Preferences) {
		for _, app := range apps {
			args.RegisteredApps = slices.DeleteFunc(args.RegisteredApps, func(a string) bool {
				return a == app
			})
		}
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
		Server:   "tcp://localhost:1883",
		User:     "",
		Password: "",
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
