// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package preferences

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// envPrefix is used to find preferences set in the environment.
const envPrefix = "GOHASSANYTHING_"

const (
	PrefServer      = "mqttserver"
	PrefUser        = "mqttuser"
	PrefPassword    = "mqttpassword"
	PrefTopicPrefix = "topicprefix"
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

//nolint:unused // these are reserved for future use
var (
	gitVersion, gitCommit, gitTreeState, buildDate string
	AppVersion                                     = gitVersion
)

type Preferences struct {
	data *koanf.Koanf
}

func findPreferences() string {
	return filepath.Join(preferencesDir, preferencesFile)
}

func (p *Preferences) GetString(key string) string {
	return p.data.String(key)
}

func (p *Preferences) GetInt(key string) int {
	return p.data.Int(key)
}

func (p *Preferences) Set(key string, value any) error {
	return p.data.Set(key, value)
}

func (p *Preferences) Keys() []string {
	return p.data.Keys()
}

func (p *Preferences) Save() error {
	data, err := p.data.Marshal(toml.Parser())
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(findPreferences(), data, 0o600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

func Load() (*Preferences, error) {
	k := koanf.New(".")
	// Load config from file.
	if err := k.Load(file.Provider(findPreferences()), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	// Merge environment variables with config.
	if err := k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("error merging file and environment: %w", err)
	}

	prefs := &Preferences{data: k}

	// Set a default server if it is not set.
	if prefs.GetString(PrefServer) == "" {
		if err := prefs.Set(PrefServer, defaultServer); err != nil {
			slog.Error("Could not set a default server", "error", err.Error())
		}
	}
	// Set a default topic prefix if it is not set.
	if prefs.GetString(PrefTopicPrefix) == "" {
		if err := prefs.Set(PrefTopicPrefix, defaultTopicPrefix); err != nil {
			slog.Error("Could not set a default topic prefix", "error", err.Error())
		}
	}

	return prefs, nil
}
