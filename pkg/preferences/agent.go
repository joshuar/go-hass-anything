// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//revive:disable:unused-receiver
package preferences

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	prefsEnvPrefix  = "GOHASSANYTHING_"
	PrefServer      = "mqtt.server"
	PrefUser        = "mqtt.user"
	PrefPassword    = "mqtt.password"
	PrefTopicPrefix = "mqtt.topicprefix"
	// defaultServer is the default MQTT broker URL.
	defaultServer = "tcp://localhost:1883"
	// defaultTopicPrefix is the default prefix that is appended to topics.
	defaultTopicPrefix = "homeassistant"
	// defaultFilePerms sets the permissions on the config file.
	defaultFilePerms = 0o600
)

var (
	// preferencesDir is the default path under which the preferences are
	// written. It can be overridden in the Save/Load functions as needed.
	preferencesDir = filepath.Join(xdg.ConfigHome, "go-hass-anything")
	// preferencesFile is the default filename used for storing the preferences
	// on disk. While it can be overridden, this is usually unnecessary.
	preferencesFile = "agent.toml"
)

// Consistent error messages.
var (
	ErrLoadPreferences     = errors.New("error loading preferences")
	ErrSavePreferences     = errors.New("error saving preferences")
	ErrValidatePreferences = errors.New("error validating preferences")
	ErrSetPreference       = errors.New("error setting preference")
)

var (
	prefsSrc = koanf.New(".")
	mu       = sync.Mutex{}
	Agent    = &AgentPreferences{}
)

// Load will retrieve the current preferences from the preference file on disk.
// If there is a problem during retrieval, an error will be returned.
var Load = func() error {
	return sync.OnceValue(func() error {
		slog.Debug("Loading preferences.", slog.String("file", filepath.Join(preferencesDir, preferencesFile)))

		// Load config file
		if err := prefsSrc.Load(file.Provider(filepath.Join(preferencesDir, preferencesFile)), toml.Parser()); err != nil {
			return fmt.Errorf("%w: %w", ErrLoadPreferences, err)
		}
		// Merge config with any environment variables.
		if err := prefsSrc.Load(env.Provider(prefsEnvPrefix, ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, prefsEnvPrefix)), "_", ".", -1)
		}), nil); err != nil {
			return fmt.Errorf("%w: %w", ErrLoadPreferences, err)
		}

		return nil
	})()
}

// Save will save the new values of the specified preferences to the existing
// preferences file. NOTE: if the preferences file does not exist, Save will
// return an error. Use New if saving preferences for the first time.
func Save() error {
	mu.Lock()
	defer mu.Unlock()

	slog.Debug("Saving preferences.", slog.String("file", filepath.Join(preferencesDir, preferencesFile)))

	if err := checkPath(preferencesDir); err != nil {
		return err
	}

	data, err := prefsSrc.Marshal(toml.Parser())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSavePreferences, err)
	}

	err = os.WriteFile(filepath.Join(preferencesDir, preferencesFile), data, defaultFilePerms)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSavePreferences, err)
	}

	return nil
}

type AgentPreferences struct{}

func (p *AgentPreferences) TopicPrefix() string {
	if prefsSrc.String(PrefTopicPrefix) == "" {
		if err := prefsSrc.Set(PrefTopicPrefix, defaultTopicPrefix); err != nil {
			slog.Warn("Could not set default value for topic prefix.", slog.Any("error", err))
		}
	}

	return prefsSrc.String(PrefTopicPrefix)
}

func (p *AgentPreferences) Server() string {
	if prefsSrc.String(PrefServer) == "" {
		if err := prefsSrc.Set(PrefServer, defaultServer); err != nil {
			slog.Warn("Could not set default value for server.", slog.Any("error", err))
		}
	}

	return prefsSrc.String(PrefServer)
}

func (p *AgentPreferences) User() string {
	return prefsSrc.String(PrefUser)
}

func (p *AgentPreferences) Password() string {
	return prefsSrc.String(PrefPassword)
}

func (p *AgentPreferences) Keys() []string {
	return []string{PrefServer, PrefUser, PrefPassword, PrefTopicPrefix}
}

func (p *AgentPreferences) GetValue(key string) (value any, found bool) {
	switch key {
	case PrefTopicPrefix:
		return p.TopicPrefix(), true
	case PrefServer:
		return p.Server(), true
	case PrefUser:
		return p.User(), true
	case PrefPassword:
		return p.Password(), true
	default:
		return nil, false
	}
}

func (p *AgentPreferences) GetDescription(key string) string {
	switch key {
	case PrefTopicPrefix:
		return "The topic prefix on which to send/receive MQTT messages."
	case PrefServer:
		return "The MQTT server formatted as a URI (e.g., tcp://localhost:1883)."
	case PrefUser:
		return "The username (when required) for connecting to MQTT."
	case PrefPassword:
		return "The password (when required) for connecting to MQTT."
	default:
		return "No description provided."
	}
}

func (p *AgentPreferences) IsSecret(key string) bool {
	return key == PrefPassword
}

func (p *AgentPreferences) SetValue(key string, value any) error {
	return prefsSrc.Set(key, value)
}
