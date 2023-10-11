// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	_ "embed"
	"errors"

	viperCustom "github.com/joshuar/go-hass-anything/pkg/config/viper"
	"github.com/spf13/viper"
)

//go:generate sh -c "printf %s $(git tag | tail -1) > VERSION"
//go:embed VERSION
var AppVersion string

//go:generate moq -out mock_configAppConfig_test.go . AppConfig
type AppConfig interface {
	// IsRegistered will retrieve the registration status of the app from its config.
	IsRegistered() bool
	// Register will set the registration status of the app.
	Register() error
	// Unregister will set the app to be unregistered.
	UnRegister() error
	// Get retrieves the value of the config key specified into the value variable passed in.
	Get(string, interface{}) error
	// Set will set the value of the config key to the specified value passed in.
	Set(string, interface{}) error
	// Delete will remove a config key from the app registration.
	Delete(string) error
}

//go:generate moq -out mock_configAgentConfig_test.go . AgentConfig
type AgentConfig interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
	Delete(string) error
}

type ConfigFileNotFoundError struct {
	Err error
}

func (e *ConfigFileNotFoundError) Error() string {
	return e.Err.Error()
}

func LoadConfig(name string) (AppConfig, error) {
	var c AppConfig
	if c, err := viperCustom.LoadViperConfig(name); err != nil && errors.Is(err, err.(viper.ConfigFileNotFoundError)) {
		return c, &ConfigFileNotFoundError{
			Err: err,
		}
	}
	return c, nil
}
