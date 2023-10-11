// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	_ "embed"

	viper "github.com/joshuar/go-hass-anything/pkg/config/viper"
)

//go:generate sh -c "printf %s $(git tag | tail -1) > VERSION"
//go:embed VERSION
var AppVersion string

//go:generate moq -out mock_configAppConfig_test.go . AppConfig
type AppConfig interface {
	IsRegistered(string) bool
	Register(string) error
	UnRegister(string) error
	Get(string, interface{}) error
	Set(string, interface{}) error
	Delete(string) error
}

//go:generate moq -out mock_configAgentConfig_test.go . AgentConfig
type AgentConfig interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
	Delete(string) error
}

func LoadConfig(name string) (AppConfig, error) {
	return viper.LoadViperConfig(name)
}
