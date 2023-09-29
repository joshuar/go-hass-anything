// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ui

//go:generate moq -out mock_Agent_test.go . Agent
type Agent interface {
	AppVersion() string
	AppName() string
	AppID() string
	Stop()
	GetConfig(string, interface{}) error
	SetConfig(string, interface{}) error
}
