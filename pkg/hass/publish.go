// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import "github.com/joshuar/go-hass-anything/pkg/mqtt"

type MQTTDevice interface {
	Configuration() []*mqtt.MQTTMsg
	States() []*mqtt.MQTTMsg
	Subscriptions() []*mqtt.MQTTSubscription
}

type MQTTClient interface {
	Publish(...*mqtt.MQTTMsg) error
	Subscribe(...*mqtt.MQTTSubscription) error
}

func PublishConfigs(device MQTTDevice, client MQTTClient) error {
	return client.Publish(device.Configuration()...)
}

func PublishState(device MQTTDevice, client MQTTClient) error {
	return client.Publish(device.States()...)
}

func Subscribe(device MQTTDevice, client MQTTClient) error {
	return client.Subscribe(device.Subscriptions()...)
}

func Unpublish(device MQTTDevice, client MQTTClient) error {
	var msgs []*mqtt.MQTTMsg
	for _, msg := range device.Configuration() {
		msgs = append(msgs, &mqtt.MQTTMsg{
			Topic:    msg.Topic,
			Message:  []byte(``),
			Retained: true,
		})
	}
	return client.Publish(msgs...)
}
