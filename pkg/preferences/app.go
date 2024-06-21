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
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type App interface {
	Name() string
}

func findAppPreferences(name string) string {
	a := strcase.ToSnake(name)

	return filepath.Join(preferencesDir, a+"_preferences.toml")
}

func LoadAppPreferences(app App) (map[string]any, error) {
	k := koanf.New(".")
	// Load config from file.
	if err := k.Load(file.Provider(findAppPreferences(app.Name())), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	prefs := make(map[string]any)

	if err := k.Unmarshal(".", &prefs); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return prefs, nil
}

func SaveAppPreferences(app App, prefs map[string]any) error {
	k := koanf.New(".")

	err := k.Load(confmap.Provider(prefs, "."), nil)
	if err != nil {
		return fmt.Errorf("could not load given preferences: %w", err)
	}

	data, err := k.Marshal(toml.Parser())
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(findAppPreferences(app.Name()), data, 0o600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}
