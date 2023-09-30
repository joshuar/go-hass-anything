// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package helpers

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type EntityConfig struct {
	Entity      *Entity
	Callback    func(MQTT.Client, MQTT.Message)
	ConfigTopic string
}

type Entity struct {
	Origin            *Origin `json:"origin,omitempty"`
	Device            *Device `json:"device,omitempty"`
	DeviceClass       string  `json:"device_class,omitempty"`
	StateTopic        string  `json:"state_topic"`
	StateClass        string  `json:"state_class,omitempty"`
	CommandTopic      string  `json:"command_topic,omitempty"`
	UnitOfMeasurement string  `json:"unit_of_measurement,omitempty"`
	ValueTemplate     string  `json:"value_template"`
	UniqueID          string  `json:"unique_id"`
	Name              string  `json:"name"`
	EntityCategory    string  `json:"entity_category,omitempty"`
}

type Device struct {
	Name          string   `json:"name"`
	Manufacturer  string   `json:"manufacturer,omitempty"`
	Model         string   `json:"model,omitempty"`
	HWVersion     string   `json:"hw_version,omitempty"`
	SWVersion     string   `json:"sw_version,omitempty"`
	URL           string   `json:"configuration_url,omitempty"`
	SuggestedArea string   `json:"suggested_area,omitempty"`
	Identifiers   []string `json:"identifiers"`
	Connections   []string `json:"connections,omitempty"`
}

type Origin struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"sw_version,omitempty"`
	URL     string `json:"support_url,omitempty"`
}
