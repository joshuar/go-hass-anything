// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"github.com/rs/zerolog/log"

	"github.com/joshuar/go-hass-anything/v3/pkg/config"
	"github.com/joshuar/go-hass-anything/v3/pkg/mqtt"
)

type MQTTDevice interface {
	Name() string
	Configuration() []*mqtt.Msg
	States() []*mqtt.Msg
	Subscriptions() []*mqtt.Subscription
}

type MQTTClient interface {
	Publish(...*mqtt.Msg) error
	Subscribe(...*mqtt.Subscription) error
}

// Register will check if the app has been registered and if not, publish the
// app configs to MQTT, which will in turn register the app with Home Assistant.
// If successfully registered, it will also record this status in the app
// configuration. If any of these actions are unsuccessful, it will return an
// error with more details. Otherwise it returns nil.
func Register(device MQTTDevice, client MQTTClient) error {
	var cfg config.AppConfig
	var err error
	if cfg, err = config.LoadConfig(device.Name()); err != nil {
		if _, ok := err.(*config.ConfigFileNotFoundError); ok {
			return err
		}
	}
	if !cfg.IsRegistered() {
		if err := PublishConfigs(device, client); err != nil {
			return err
		} else {
			if err = cfg.Register(); err != nil {
				return err
			}
		}
		log.Debug().Str("appName", device.Name()).Msg("App registered.")
	}
	return nil
}

// UnRegister performs two actions. First, it removes the app config topics from
// MQTT, effectively removing the app from Home Assistant. Second, it updates
// the app config to indicate it is unregistered. It will return an error if
// either action fails, otherwise it will return nil.
func UnRegister(device MQTTDevice, client MQTTClient) error {
	var cfg config.AppConfig
	var err error
	if cfg, err = config.LoadConfig(device.Name()); err != nil {
		if _, ok := err.(*config.ConfigFileNotFoundError); ok {
			return err
		}
	}
	if err := Unpublish(device, client); err != nil {
		return err
	}
	if err := cfg.UnRegister(); err != nil {
		return err
	}
	log.Debug().Str("appName", device.Name()).Msg("App unregistered.")
	return nil
}

func PublishConfigs(device MQTTDevice, client MQTTClient) error {
	return client.Publish(device.Configuration()...)
}

func PublishState(device MQTTDevice, client MQTTClient) error {
	log.Debug().Str("appName", device.Name()).Msg("Publishing messages.")
	return client.Publish(device.States()...)
}

func Subscribe(device MQTTDevice, client MQTTClient) error {
	return client.Subscribe(device.Subscriptions()...)
}

func Unpublish(device MQTTDevice, client MQTTClient) error {
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
