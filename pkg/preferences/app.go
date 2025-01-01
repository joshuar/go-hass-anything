// Copyright 2024 vscode.
// SPDX-License-Identifier: MIT

package preferences

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/pelletier/go-toml"

	"github.com/joshuar/go-hass-anything/v12/pkg/validation"
)

// App represents an App from the point of the preferences package. An
// app has a set of default preferences returned by the DefaultPreferences
// method and an ID that uniquely identifies the app (and its preferences
// on disk).
type App[T any] interface {
	PreferencesID() string
	DefaultPreferences() T
}

var (
	ErrSaveAppPrefs = errors.New("error saving app preferences")
	ErrLoadAppPrefs = errors.New("error loading app preferences")
)

// LoadApp reads the given apps's preferences from file.
func LoadApp[T any](app App[T]) (*T, error) {
	prefsKey := "app." + app.PreferencesID()
	// Load default app prefs.
	prefs := app.DefaultPreferences()

	if prefsSrc.Get(prefsKey) == nil {
		slog.Debug("Using default preferences for app.",
			slog.String("app", app.PreferencesID()))
		// Save the default preferences to the preferences source.
		if err := SaveApp(app, prefs); err != nil {
			return &prefs, fmt.Errorf("%w: %w", ErrLoadAppPrefs, err)
		}

		return &prefs, nil
	}

	// Unmarshal the existing prefs into the prefs type, overwriting any
	// defaults.
	if err := prefsSrc.Unmarshal(prefsKey, &prefs); err != nil {
		return &prefs, fmt.Errorf("%w: %w", ErrLoadAppPrefs, err)
	}

	// If the preferences are invalid, warn and use defaults.
	if err := validation.Validate.Struct(prefs); err != nil {
		slog.Warn("App preferences are invalid, using defaults.",
			slog.String("app", app.PreferencesID()),
			slog.String("problems", validation.ParseValidationErrors(err)))

		prefs = app.DefaultPreferences()

		return &prefs, nil
	}

	// Return preferences.
	return &prefs, nil
}

// SaveApp saves the given app's preferences to file.
func SaveApp[T any](app App[T], prefs T) error {
	// We can't define the structure for every possible app beforehand, so
	// use map[string]any as the structure for saving.
	prefsMaps := make(map[string]any)

	// Marshal the app's prefs object into bytes, using the toml tag
	// structure.
	data, err := toml.Marshal(&prefs)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSaveAppPrefs, err)
	}
	// Unmarshal back into a map[string]any that we can save into the preferences
	// file.
	if err := toml.Unmarshal(data, &prefsMaps); err != nil {
		return fmt.Errorf("%w: %w", ErrSaveAppPrefs, err)
	}

	// Merge the app preferences into the preferences file.
	if err := prefsSrc.Set("app."+app.PreferencesID(), prefsMaps); err != nil {
		return fmt.Errorf("%w: %w", ErrSaveAppPrefs, err)
	}

	return nil
	// Save the preferences.
	// return Save()
}
