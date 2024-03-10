// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"errors"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joshuar/go-hass-anything/v6/pkg/mqtt"
	"github.com/joshuar/go-hass-anything/v6/pkg/preferences"
	"github.com/rs/zerolog/log"
)

var (
	ErrRegistrationFailed = errors.New("registration failed")
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
func Register(client MQTTClient, device ...MQTTDevice) error {
	for _, dev := range device {
		log.Debug().Str("app", dev.Name()).Msg("Registering app.")
		if err := client.Publish(dev.Configuration()...); err != nil {
			return ErrRegistrationFailed
		}
	}
	prefs, err := preferences.LoadPreferences()
	if err != nil {
		return err
	}
	HAStatusSub := &mqtt.Subscription{
		Topic: prefs.GetTopicPrefix() + "/status",
		Callback: func(c MQTT.Client, m MQTT.Message) {
			switch msg := string(m.Payload()); msg {
			case "online":
				log.Debug().Msg("Home Assistant Online.")
				for _, dev := range device {
					log.Debug().Str("app", dev.Name()).Msg("Re-registering app.")
					if err := client.Publish(dev.Configuration()...); err != nil {
						log.Warn().Err(err).Str("app", dev.Name()).Msg("Could not re-register.")
						return
					}
					if err := PublishState(dev, client); err != nil {
						log.Warn().Err(err).Str("app", dev.Name()).Msg("Could not publish state.")
					}

				}
			case "offline":
				log.Debug().Msg("Home Assistant Offline.")
			}
		},
	}
	if err := client.Subscribe(HAStatusSub); err != nil {
		return errors.Join(ErrRegistrationFailed, err)
	}
	return nil
}

// UnRegister removes the app config topics from MQTT, effectively removing the
// app from Home Assistant. It will return an error if it fails, otherwise it
// will return nil.
func UnRegister(client MQTTClient, device ...MQTTDevice) error {
	var msgs []*mqtt.Msg
	for _, dev := range device {
		log.Debug().Str("app", dev.Name()).Msg("Unregistering app.")
		for _, msg := range dev.Configuration() {
			msgs = append(msgs, mqtt.NewMsg(msg.Topic, []byte(``)))
		}
		if err := client.Publish(msgs...); err != nil {
			return err
		}
	}
	return nil
}

func PublishState(device MQTTDevice, client MQTTClient) error {
	return client.Publish(device.States()...)
}

func Subscribe(device MQTTDevice, client MQTTClient) error {
	return client.Subscribe(device.Subscriptions()...)
}
