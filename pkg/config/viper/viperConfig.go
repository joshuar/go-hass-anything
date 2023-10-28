// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	PrefAppRegistered = "app.registered"
)

var (
	configBasePath = filepath.Join(os.Getenv("HOME"), ".config", "go-hass-anything")
)

type ViperConfig struct {
	store *viper.Viper
	path  string
}

func (c *ViperConfig) Get(key string, value interface{}) error {
	switch v := value.(type) {
	case *string:
		*v = c.store.GetString(key)
		if *v == "" {
			return fmt.Errorf("key %s not set", key)
		}
	case *bool:
		*v = c.store.GetBool(key)
	default:
		return errors.New("unsupported config value")
	}
	return nil
}

func (c *ViperConfig) Set(key string, value interface{}) error {
	c.store.Set(key, value)
	if err := c.store.WriteConfigAs(filepath.Join(c.path, "config.toml")); err != nil {
		log.Error().Err(err).Msg("Problem writing config file.")
	}
	return nil
}

func (c *ViperConfig) Delete(key string) error {
	return nil
}

func (c *ViperConfig) IsRegistered() bool {
	return c.store.GetBool(PrefAppRegistered)
}

func (c *ViperConfig) Register() error {
	return c.Set(PrefAppRegistered, true)
}

func (c *ViperConfig) UnRegister() error {
	return c.Set(PrefAppRegistered, false)
}

func LoadViperConfig(name string) (*ViperConfig, error) {
	c := &ViperConfig{
		store: viper.New(),
	}
	if name != "" {
		c.path = filepath.Join(configBasePath, name)
	} else {
		c.path = configBasePath
	}
	c.store.SetConfigName("config")
	c.store.SetConfigType("toml")
	c.store.AddConfigPath(c.path)
	if err := c.store.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			createDir(c.path)
			return c, err
		} else {
			return nil, err
		}
	}
	log.Debug().Msgf("Config file %s", c.store.ConfigFileUsed())
	return c, nil
}

func createDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Debug().Err(err).Msgf("Failed to create new config directory %s.", path)
		} else {
			log.Debug().Msgf("Created new config directory %s.", path)
		}
	}
}
