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

type prefs interface {
	GetMQTTServer() string
	GetMQTTUser() string
	GetMQTTPassword() string
	GetTopicPrefix() string
}

type Device interface {
	Name() string
	Configuration() []*Msg
	States() []*Msg
	Subscriptions() []*Subscription
}

// Subscription represents a listener on a specific Topic, that will pass any
// messages sent to that topic to the Callback function.
type Subscription struct {
	Callback func(MQTT.Client, MQTT.Message)
	Topic    string
	QOS      byte
	Retained bool
}

// Client is the connection to the MQTT broker.
type Client struct {
	conn    MQTT.Client
	options *MQTT.ClientOptions
}

// Publish will send the list of messages it is passed to the broker that the
// client is connected to. Any errors in publihsing will be returned.
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

// Subscribe will parse the list of subscriptions and listen on their topics,
// passing any received messages to their callback functions. Any error in
// setting up a subscription will be returned.
func (c *Client) Subscribe(subs ...*Subscription) error {
	g, _ := errgroup.WithContext(context.TODO())
	msgCh := make(chan *Subscription, len(subs))

	for i := 0; i < len(subs); i++ {
		msgCh <- subs[i]
	}

	g.Go(func() error {
		for sub := range msgCh {
			log.Trace().Str("topic", sub.Topic).Bool("retain", sub.Retained).Msg("Subscribing to topic.")
			if token := c.conn.Subscribe(sub.Topic, sub.QOS, sub.Callback); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
		return nil
	})

	close(msgCh)
	return g.Wait()
}

func (c *Client) Unpublish(msgs ...*Msg) error {
	for _, msg := range msgs {
		msgs = append(msgs, NewMsg(msg.Topic, []byte(``)))
	}
	if err := c.Publish(msgs...); err != nil {
		return err
	}
	return nil
}

// Connect will establish a new connection to the MQTT service with a generated configuration
func (c *Client) connect(ctx context.Context) error {
	c.conn = MQTT.NewClient(c.options)

	connect := func() error {
		if token := c.conn.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}
		return nil
	}
	err := backoff.Retry(connect, backoff.WithContext(backoff.NewExponentialBackOff(), ctx))
	if err != nil {
		return err
	}

	log.Debug().Msg("Connected to MQTT server.")
	return nil
}

func NewClient(ctx context.Context, prefs prefs, devices ...Device) (*Client, error) {
	statusTopic := prefs.GetTopicPrefix() + "/status"
	onConnectCallback := genOnConnectHandler(statusTopic, devices...)
	connOpts := genConnOpts(prefs, onConnectCallback)
	c := &Client{
		options: connOpts,
	}
	if err := c.connect(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func genConnOpts(prefs prefs, callback MQTT.OnConnectHandler) *MQTT.ClientOptions {
	hostname, _ := os.Hostname()
	clientid := hostname + strconv.Itoa(time.Now().Second())

	connOpts := MQTT.NewClientOptions().
		AddBroker(prefs.GetMQTTServer()).
		SetClientID(clientid).
		SetCleanSession(true).
		SetKeepAlive(10 * time.Second).
		SetAutoReconnect(true).
		SetOnConnectHandler(callback)

	if prefs.GetMQTTUser() != "" {
		connOpts.SetUsername(prefs.GetMQTTUser())
		if prefs.GetMQTTPassword() != "" {
			connOpts.SetPassword(prefs.GetMQTTPassword())
		}
	}

	return connOpts
}

func genOnConnectHandler(statustopic string, devices ...Device) MQTT.OnConnectHandler {
	var configs []*Msg
	var subscriptions []*Subscription

	for _, d := range devices {
		configs = append(configs, d.Configuration()...)
		subscriptions = append(subscriptions, d.Subscriptions()...)
	}
	HAStatusSub := &Subscription{
		Topic: statustopic,
		Callback: func(c MQTT.Client, m MQTT.Message) {
			switch msg := string(m.Payload()); msg {
			case "online":
				log.Debug().Msg("Home Assistant Online.")
				for _, device := range devices {
					publish(c, device.Configuration()...)
					publish(c, device.States()...)
					subscribe(c, device.Subscriptions()...)
				}
			case "offline":
				log.Debug().Msg("Home Assistant Offline.")
			}
		},
	}
	subscriptions = append(subscriptions, HAStatusSub)
	redoFunc := func(c MQTT.Client) {
		publish(c, configs...)
		subscribe(c, subscriptions...)
	}
	return redoFunc
}

func publish(c MQTT.Client, msgs ...*Msg) {
	for _, msg := range msgs {
		log.Trace().Str("topic", msg.Topic).Bool("retain", msg.Retained).Msg("Publishing message.")
		if token := c.Publish(msg.Topic, msg.QOS, msg.Retained, []byte(msg.Message)); token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Str("topic", msg.Topic).Msg("Failed to publish message.")
		}
	}
}

func subscribe(c MQTT.Client, subs ...*Subscription) {
	for _, sub := range subs {
		log.Trace().Str("topic", sub.Topic).Bool("retain", sub.Retained).Msg("Subscribing to topic.")
		if token := c.Subscribe(sub.Topic, sub.QOS, sub.Callback); token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Str("topic", sub.Topic).Msg("Failed to subscribe to topic.")
		}
	}
}

// Msg represents a message that can be sent or received on the MQTT bus.
type Msg struct {
	Topic    string
	Message  json.RawMessage
	QOS      byte
	Retained bool
}

// Retain sets the Retained status of a Msg to true, ensuring that it will be
// retained on the MQTT bus when sent.
func (m *Msg) Retain() *Msg {
	m.Retained = true
	return m
}

// NewMsg is a convenience function to create a new Msg with a given topic and
// message body. The returned Msg can be further customised directly for
// specifying retention and QoS parameters.
func NewMsg(topic string, msg json.RawMessage) *Msg {
	return &Msg{
		Topic:   topic,
		Message: msg,
	}
}
