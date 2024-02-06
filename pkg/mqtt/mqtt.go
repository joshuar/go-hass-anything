// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const (
	DiscoveryPrefix = "homeassistant"
	DefaultServer   = "localhost:1883"
)

type prefs interface {
	MQTTServer() string
	MQTTUser() string
	MQTTPassword() string
}

type Msg struct {
	Topic    string
	Message  json.RawMessage
	QOS      byte
	Retained bool
}

type Subscription struct {
	Callback func(MQTT.Client, MQTT.Message)
	Topic    string
	QOS      byte
	Retained bool
}

type Client struct {
	conn MQTT.Client
}

func (c *Client) Publish(msgs ...*Msg) error {
	g, _ := errgroup.WithContext(context.TODO())
	msgCh := make(chan *Msg, len(msgs))

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
				}
				i++
			} else {
				log.Debug().Msg("Not connected.")
			}
		}
		log.Trace().Int("msgCount", i).Msg("Finished publishing messages.")
		return nil
	})

	close(msgCh)
	return g.Wait()
}

func (c *Client) Subscribe(subs ...*Subscription) error {
	g, _ := errgroup.WithContext(context.TODO())
	msgCh := make(chan *Subscription, len(subs))

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

// NewMQTTClient will establish a new connection to the MQTT service, using the
// configuration found under the path specified with prefsPath.
func NewMQTTClient(preferences prefs) (*Client, error) {
	hostname, _ := os.Hostname()
	clientid := hostname + strconv.Itoa(time.Now().Second())

	connOpts := MQTT.NewClientOptions().AddBroker(preferences.MQTTServer()).SetClientID(clientid).SetCleanSession(true)
	if preferences.MQTTUser() != "" {
		connOpts.SetUsername(preferences.MQTTUser())
		if preferences.MQTTPassword() != "" {
			connOpts.SetPassword(preferences.MQTTPassword())
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

	log.Debug().Msgf("Connected to MQTT server %s.", preferences.MQTTServer())
	conf := &Client{
		conn: client,
	}

	return conf, nil
}
