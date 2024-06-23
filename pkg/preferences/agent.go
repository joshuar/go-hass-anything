// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//revive:disable:unused-receiver
package preferences

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

const (
	PrefServer      = "mqtt.server"
	PrefUser        = "mqtt.user"
	PrefPassword    = "mqtt.password"
	PrefTopicPrefix = "mqtt.topicprefix"
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

type AgentPreferences struct {
	MQTTServer      *Preference `toml:"mqtt.server"`
	MQTTUser        *Preference `toml:"mqtt.user,omitempty"`
	MQTTPassword    *Preference `toml:"mqtt.password,omitempty"`
	MQTTTopicPrefix *Preference `toml:"mqtt.topicprefix"`
}

func (p *AgentPreferences) TopicPrefix() string {
	if pref, ok := p.MQTTTopicPrefix.Value.(string); ok {
		return pref
	}

	return ""
}

func (p *AgentPreferences) Server() string {
	if pref, ok := p.MQTTServer.Value.(string); ok {
		return pref
	}

	return ""
}

func (p *AgentPreferences) User() string {
	if pref, ok := p.MQTTUser.Value.(string); ok {
		return pref
	}

	return ""
}

func (p *AgentPreferences) Password() string {
	if pref, ok := p.MQTTPassword.Value.(string); ok {
		return pref
	}

	return ""
}

// findAgentPreferences returns the file path to the file for the agent
// preferences.
func findAgentPreferences() string {
	return filepath.Join(preferencesDir, preferencesFile)
}

func (p *AgentPreferences) Keys() []string {
	return []string{PrefServer, PrefUser, PrefPassword, PrefTopicPrefix}
}

func (p *AgentPreferences) GetValue(key string) (value any, found bool) {
	pref := p.getPref(key)
	if pref == nil {
		return nil, false
	}

	return p.getPref(key).Value, true
}

func (p *AgentPreferences) GetDescription(key string) string {
	return p.getPref(key).Description
}

func (p *AgentPreferences) IsSecret(key string) bool {
	return p.getPref(key).Secret
}

func (p *AgentPreferences) SetValue(key string, value any) error {
	switch key {
	case PrefServer:
		p.MQTTServer.Value = value
	case PrefUser:
		p.MQTTUser.Value = value
	case PrefPassword:
		p.MQTTPassword.Value = value
	case PrefTopicPrefix:
		p.MQTTTopicPrefix.Value = value
	default:
		return ErrUnknownPref
	}

	return nil
}

//nolint:mnd
func Save(prefs any) error {
	data, err := toml.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(findAgentPreferences(), data, 0o600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

//nolint:exhaustruct
func Load() (*AgentPreferences, error) {
	if err := checkPath(preferencesDir); err != nil {
		return nil, fmt.Errorf("could not create new preferences directory: %w", err)
	}

	data, err := os.ReadFile(findAgentPreferences())
	if err != nil {
		return defaultAgentPrefs(), fmt.Errorf("could not read agent preferences file: %w, using defaults", err)
	}

	prefs := &AgentPreferences{}

	if err := toml.Unmarshal(data, prefs); err != nil {
		return defaultAgentPrefs(), fmt.Errorf("could not parse agent preferences file: %w, using defaults", err)
	}

	return prefs, nil
}

func (p *AgentPreferences) getPref(key string) *Preference {
	switch key {
	case PrefServer:
		return p.MQTTServer
	case PrefUser:
		return p.MQTTUser
	case PrefPassword:
		return p.MQTTPassword
	case PrefTopicPrefix:
		return p.MQTTTopicPrefix
	default:
		return nil
	}
}

//nolint:exhaustruct
func defaultAgentPrefs() *AgentPreferences {
	return &AgentPreferences{
		MQTTServer: &Preference{
			Value:       defaultServer,
			Description: "The MQTT server (in tcp://some.host:port format).",
		},
		MQTTTopicPrefix: &Preference{
			Value:       defaultTopicPrefix,
			Description: "The MQTT topic prefix.",
		},
		MQTTUser: &Preference{
			Value:       "",
			Description: "The username required to authenticate with the MQTT server.",
		},
		MQTTPassword: &Preference{
			Value:       "",
			Description: "The password required to authenticate with the MQTT server.",
			Secret:      true,
		},
	}
}
