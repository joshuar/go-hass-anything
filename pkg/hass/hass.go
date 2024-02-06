// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"github.com/joshuar/go-hass-anything/v4/pkg/mqtt"
)

type MQTTDevice interface {
	Name() string
	Configuration() []*mqtt.Msg
	States() []*mqtt.Msg
	Subscriptions() []*mqtt.Subscription
}

type MQTTClient interface {
	Publish(msgs ...*mqtt.Msg) error
	Subscribe(msgs ...*mqtt.Subscription) error
}

// Register  publish the app configs to MQTT, which will in turn register the
// app with Home Assistant. If  unsuccessful, it will return an error with more
// details. Otherwise it returns nil.
func Register(device MQTTDevice, client MQTTClient) error {
	return client.Publish(device.Configuration()...)
}

// UnRegister removes the app config topics from MQTT, effectively removing the
// app from Home Assistant. It will return an error if it fails, otherwise it
// will return nil.
func UnRegister(device MQTTDevice, client MQTTClient) error {
	var msgs []*mqtt.Msg
	for _, msg := range device.Configuration() {
		msgs = append(msgs, &mqtt.Msg{
			Topic:    msg.Topic,
			Message:  []byte(``),
			Retained: true,
		})
	}
	return client.Publish(msgs...)
}

func PublishState(device MQTTDevice, client MQTTClient) error {
	return client.Publish(device.States()...)
}

func Subscribe(device MQTTDevice, client MQTTClient) error {
	return client.Subscribe(device.Subscriptions()...)
}
