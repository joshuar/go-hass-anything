// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package hass

import (
	"github.com/joshuar/go-hass-anything/pkg/config"
	"github.com/joshuar/go-hass-anything/pkg/mqtt"
	"github.com/rs/zerolog/log"
)

type MQTTDevice interface {
	Name() string
	Configuration() []*mqtt.MQTTMsg
	States() []*mqtt.MQTTMsg
	Subscriptions() []*mqtt.MQTTSubscription
}

type MQTTClient interface {
	Publish(...*mqtt.MQTTMsg) error
	Subscribe(...*mqtt.MQTTSubscription) error
}

func Register(device MQTTDevice, client MQTTClient) error {
	var cfg config.AppConfig
	var err error
	if cfg, err = config.LoadConfig(device.Name()); err != nil {
		return err
	}
	if !cfg.IsRegistered(device.Name()) {
		if err := PublishConfigs(device, client); err != nil {
			return err
		} else {
			if err = cfg.Register(device.Name()); err != nil {
				return err
			}
		}
		log.Debug().Str("appName", device.Name()).Msg("App registered.")
	}
	return nil
}

func UnRegister(device MQTTDevice, client MQTTClient) error {
	var cfg config.AppConfig
	var err error
	if cfg, err = config.LoadConfig(device.Name()); err != nil {
		return err
	}
	if err := Unpublish(device, client); err != nil {
		return err
	}
	if err := cfg.UnRegister(device.Name()); err != nil {
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
