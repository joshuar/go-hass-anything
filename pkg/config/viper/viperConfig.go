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

	"github.com/joshuar/go-hass-anything/pkg/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

func (c *ViperConfig) IsRegistered(name string) bool {
	return c.store.GetBool(config.PrefAppRegistered)
}

func (c *ViperConfig) Register(name string) error {
	return c.Set(config.PrefAppRegistered, true)
}

func (c *ViperConfig) UnRegister(name string) error {
	return c.Set(config.PrefAppRegistered, false)
}

func Load(name string) (*ViperConfig, error) {
	c := &ViperConfig{
		store: viper.New(),
	}
	if name != "" {
		c.path = filepath.Join(config.ConfigBasePath, name)
	} else {
		c.path = config.ConfigBasePath
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
