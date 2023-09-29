// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

//go:generate moq -out mock_configApp_test.go . App
type App interface {
	IsRegistered(string) bool
	Register(string) error
	UnRegister(string) error
	Get(string, interface{}) error
	Set(string, interface{}) error
	Delete(string) error
}

//go:generate moq -out mock_configAgent_test.go . Agent
type Agent interface {
	Get(string, interface{}) error
	Set(string, interface{}) error
	Delete(string) error
}
