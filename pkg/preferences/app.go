// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package preferences

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// AppPreferences holds a given app's preferences as a generic map[string]any.
type AppPreferences struct {
	Prefs map[string]any `toml:"preferences"`
}

// AppPref is a functional type for applying AppPreferences
type AppPref func(*AppPreferences)

// SetAppPreferences will assign the given preferences map to the AppPreferences
// struct.
func SetAppPreferences(prefs map[string]any) AppPref {
	return func(args *AppPreferences) {
		args.Prefs = prefs
	}
}

// SaveAppPreferences will save the given preferences for the given app to a
// file on disk. Any existing preferences will be preserved. If there are no
// existing preferences, a new preferences file will be created.
func SaveAppPreferences(app string, setters ...AppPref) error {
	file := filepath.Join(preferencesDir, app+"-preferences.toml")
	if err := checkPath(preferencesDir); err != nil {
		return err
	}

	prefs, err := LoadAppPreferences(app)
	if err != nil {
		return err
	}

	for _, setter := range setters {
		setter(prefs)
	}

	return write(prefs, file)
}

// LoadAppPreferences will load the given app preferences from file. If there
// are no existing preferences, it will return an os.IsNotExist error and a
// default (empty) AppPreferences.
func LoadAppPreferences(app string) (*AppPreferences, error) {
	file := filepath.Join(preferencesDir, app+"-preferences.toml")

	p := DefaultAppPreferences()
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

// DefaultAppPreferences returns an empty AppPreferences struct ready for use by
// apps.
func DefaultAppPreferences() *AppPreferences {
	return &AppPreferences{
		Prefs: make(map[string]any),
	}
}
