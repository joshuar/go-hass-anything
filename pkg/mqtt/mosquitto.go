// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joshuar/go-hass-anything/pkg/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	DiscoveryPrefix = "homeassistant"
	DefaultServer   = "localhost:1883"
)

//go:generate moq -out mock_AgentConfig_test.go . AgentConfig
type AgentConfig interface {
	GetConfig(string, interface{}) error
}

type mqttConfig struct {
	server   string
	user     string
	password string
}

type MQTTMsg struct {
	Topic    string
	Message  json.RawMessage
	QOS      byte
	Retained bool
}

type MQTTSubscription struct {
	Callback func(MQTT.Client, MQTT.Message)
	Topic    string
	QOS      byte
	Retained bool
}

type MQTTClient struct {
	conn MQTT.Client
}

func (c *MQTTClient) Publish(msgs ...*MQTTMsg) error {
	g, _ := errgroup.WithContext(context.TODO())
	msgCh := make(chan *MQTTMsg, len(msgs))

	for i := 0; i < len(msgs); i++ {
		msgCh <- msgs[i]
	}

	g.Go(func() error {
		var i int
		for msg := range msgCh {
			log.Trace().Str("topic", msg.Topic).Bool("retain", msg.Retained).Msg("Publishing message.")
			if c.conn.IsConnected() {
				if token := c.conn.Publish(msg.Topic, msg.QOS, msg.Retained, []byte(msg.Message)); token.Wait() && token.Error() != nil {
					return token.Error()
				} else {
					i++
				}
			} else {
				log.Debug().Msg("Not connected.")
			}
		}
		log.Debug().Int("msgCount", i).Msg("Finished publishing messages.")
		return nil
	})

	close(msgCh)
	return g.Wait()
}

func (c *MQTTClient) Subscribe(subs ...*MQTTSubscription) error {
	g, _ := errgroup.WithContext(context.TODO())
	msgCh := make(chan *MQTTSubscription, len(subs))

	for i := 0; i < len(subs); i++ {
		msgCh <- subs[i]
	}

	g.Go(func() error {
		for sub := range msgCh {
			log.Trace().Str("topic", sub.Topic).Bool("retain", sub.Retained).Msg("Adding subscription.")
			if token := c.conn.Subscribe(sub.Topic, sub.QOS, sub.Callback); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
		return nil
	})

	close(msgCh)
	return g.Wait()
}

func NewMQTTClient(cfg AgentConfig) (*MQTTClient, error) {
	hostname, _ := os.Hostname()
	clientid := hostname + strconv.Itoa(time.Now().Second())

	c := &mqttConfig{}
	if err := cfg.GetConfig(config.PrefMQTTServer, &c.server); err != nil {
		return nil, err
	}
	if c.server == "" || c.server == "NOTSET" {
		return nil, errors.New("invalid server value")
	}

	connOpts := MQTT.NewClientOptions().AddBroker(c.server).SetClientID(clientid).SetCleanSession(true)
	if c.user != "" {
		connOpts.SetUsername(c.user)
		if c.password != "" {
			connOpts.SetPassword(c.password)
		}
	}

	client := MQTT.NewClient(connOpts)

	connect := func() error {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		return nil
	}
	err := backoff.Retry(connect, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("Connected to MQTT server %s.", c.server)
	conf := &MQTTClient{
		conn: client,
	}

	return conf, nil
}
