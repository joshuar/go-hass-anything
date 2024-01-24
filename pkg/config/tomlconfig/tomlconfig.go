// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package tomlconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
)

const (
	PrefAppRegistered = "registered"
)

var configBasePath = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")

type tomlConfig struct {
	file string
}

func (c *tomlConfig) read() (map[string]any, error) {
	items := make(map[string]any)
	b, err := os.ReadFile(c.file)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(b, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (c *tomlConfig) write(items map[string]any) error {
	b, err := toml.Marshal(items)
	if err != nil {
		return err
	}
	err = os.WriteFile(c.file, b, 0o600)
	if err != nil {
		return err
	}
	return nil
}

func (c *tomlConfig) Get(key string, value any) error {
	items, err := c.read()
	if err != nil {
		return err
	}
	item, ok := items[key]
	if !ok {
		return fmt.Errorf("key %s not set", key)
	}
	switch v := value.(type) {
	case *string:
		*v = item.(string)
	case *bool:
		*v = item.(bool)
	default:
		return errors.New("unsupported config value")
	}
	return nil
}

func (c *tomlConfig) Set(key string, value any) error {
	items, err := c.read()
	if err != nil {
		return err
	}
	items[key] = value
	return c.write(items)
}

func (c *tomlConfig) Delete(key string) error {
	items, err := c.read()
	if err != nil {
		return err
	}
	delete(items, key)
	return c.write(items)
}

func (c *tomlConfig) IsRegistered() bool {
	var r bool
	err := c.Get(PrefAppRegistered, &r)
	if err != nil {
		log.Warn().Err(err).Msg("Could not get registration status.")
		return false
	}
	return r
}

func (c *tomlConfig) Register() error {
	return c.Set(PrefAppRegistered, true)
}

func (c *tomlConfig) UnRegister() error {
	return c.Set(PrefAppRegistered, false)
}

func newTomlConfig(path string) *tomlConfig {
	c := &tomlConfig{
		file: filepath.Join(path, "config.toml"),
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Debug().Err(err).Msgf("Failed to create new config directory %s.", path)
		} else {
			log.Debug().Msgf("Created new config directory %s.", path)
			i := make(map[string]any)
			i[PrefAppRegistered] = false
			err := c.write(i)
			if err != nil {
				log.Debug().Err(err).Msg("Could not set initial registration state.")
			}
		}
	}
	return c
}

func LoadTOMLConfig(name string) *tomlConfig {
	var path string
	if name != "" {
		path = filepath.Join(configBasePath, name)
	} else {
		path = configBasePath
	}
	return newTomlConfig(path)
}
