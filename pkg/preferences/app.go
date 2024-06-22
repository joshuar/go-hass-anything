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
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func findAppPreferences(name string) string {
	a := strcase.ToSnake(name)

	return filepath.Join(preferencesDir, a+"_preferences.toml")
}

func (p *Preferences) SaveApp(app string) error {
	data, err := p.data.Marshal(toml.Parser())
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(findAppPreferences(app), data, 0o600); err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}

	return nil
}

func LoadApp(app string) (*Preferences, error) {
	k := koanf.New(".")
	// Load config from file.
	if err := k.Load(file.Provider(findAppPreferences(app)), toml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config file: %w", err)
	}

	prefs := &Preferences{data: k}

	return prefs, nil
}
