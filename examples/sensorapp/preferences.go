// Copyright 2024 vscode.
// SPDX-License-Identifier: MIT

//revive:disable:unused-receiver
package sensorapp

import (
	"errors"
	"fmt"
	"net/url"
)

// Prefs represent our app preferences.
type Prefs struct {
	URL string `toml:"weather_url"`
}

// Keys returns an array of the prefererence keys.
func (p *Prefs) Keys() []string {
	return []string{weatherURLpref}
}

// GetValue will fetch a particular preference.
func (p *Prefs) GetValue(key string) (value any, found bool) {
	switch key {
	case weatherURLpref:
		return p.URL, true
	default:
		return nil, false
	}
}

// SetValue will save a particular preference.
func (p *Prefs) SetValue(key string, value any) error {
	switch key {
	case weatherURLpref:
		str, ok := value.(string)
		if !ok {
			return errors.New("invalid weather service url: not a string")
		}

		if _, err := url.Parse(str); err != nil {
			return fmt.Errorf("invalid weather service url: %w", err)
		}

		p.URL = str

		return nil
	default:
		return fmt.Errorf("unknown preference: %s", key)
	}
}

// GetDescription returns a textual description of a preference. Called when
// displaying preferences in the agent TUI during configuration.
func (p *Prefs) GetDescription(key string) string {
	switch key {
	case weatherURLpref:
		return "The weather service to use for fetching weather (default:" + weatherURL + ")."
	default:
		return "No description provided."
	}
}
