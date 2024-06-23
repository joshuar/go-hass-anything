// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package preferences

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/pelletier/go-toml"
)

func findAppPreferences(name string) string {
	a := strcase.ToSnake(name)

	return filepath.Join(preferencesDir, a+"_preferences.toml")
}

// AppPreferences is a structure that can be used to represent app preferences.
// As app preferences vary between apps, a map of Preference values is used.
type AppPreferences map[string]*Preference

// GetValue returns the value of the preference with the given key name and a
// bool to indicate whether it was found or not. If the preference was not
// found, the value will be nil.
func (p AppPreferences) GetValue(key string) (value any, found bool) {
	value, found = p[key]
	if !found {
		return nil, false
	}

	return value, true
}

// SetValue sets the preference with the given key name to the given value. It
// currently returns a nil error but may in the future return a non-nil error if
// the preference could not be set.
func (p AppPreferences) SetValue(key string, value any) error {
	p[key].Value = value

	return nil
}

// GetDescription returns the description of the preference with the given key
// name.
func (p AppPreferences) GetDescription(key string) string {
	return p[key].Description
}

// IsSecret returns a boolean to indicate whether the preference with the given
// key name should be masked or obscured when displaying to the user.
func (p AppPreferences) IsSecret(key string) bool {
	return p[key].Secret
}

// Keys returns all key names for all known preferences.
func (p AppPreferences) Keys() []string {
	keys := make([]string, 0, len(p))
	for key := range p {
		keys = append(keys, key)
	}

	return keys
}

// SaveApp will save the given preferences for the given app to file. If the
// preferences cannot be saved, a non-nil error will be returned.
//
//nolint:mnd
func SaveApp(app string, prefs AppPreferences) error {
	data, err := toml.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("could not marshal app preferences: %w", err)
	}

	if err := os.WriteFile(findAppPreferences(app), data, 0o600); err != nil {
		return fmt.Errorf("could not write app preferences: %w", err)
	}

	return nil
}

// LoadApp will load the given app preferences from file. If the preferences
// cannot be loaded, a non-nil error will be returned. Apps should take special
// care to handle os.ErrNotExists verses other returned errors. In the former case, it
// would be wise to treat as not an error and revert to using default
// preferences.
func LoadApp(app string) (AppPreferences, error) {
	// Load config from file.
	data, err := os.ReadFile(findAppPreferences(app))
	if err != nil {
		return nil, fmt.Errorf("could not read app preferences file: %w", err)
	}

	prefs := make(AppPreferences)

	if err := toml.Unmarshal(data, &prefs); err != nil {
		return nil, fmt.Errorf("could not unmarshal app preferences: %w", err)
	}

	return prefs, nil
}
