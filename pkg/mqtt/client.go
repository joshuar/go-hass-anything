// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

import (
	"context"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/rs/zerolog/log"
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
	Callback func(p *paho.Publish)
	Topic    string
}

// Client is the connection to the MQTT broker.
type Client struct {
	conn *autopaho.ConnectionManager
	mu   sync.Mutex
}

// Publish will send the list of messages it is passed to the broker that the
// client is connected to. Any errors in publihsing will be returned.
func (c *Client) Publish(msgs ...*Msg) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := publish(c.conn, msgs...)
	return err
}

func (c *Client) Unpublish(msgs ...*Msg) error {
	var newMsgs []*Msg
	for _, msg := range msgs {
		newMsgs = append(newMsgs, NewMsg(msg.Topic, []byte(``)))
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	err := publish(c.conn, newMsgs...)
	return err
}

func NewClient(ctx context.Context, prefs prefs, subscriptions []*Subscription) (*Client, error) {
	var subOpts []paho.SubscribeOptions
	router := paho.NewStandardRouter()

	for _, s := range subscriptions {
		log.Trace().Str("topic", s.Topic).Msg("Adding subscription for topic.")
		subOpts = append(subOpts, paho.SubscribeOptions{Topic: s.Topic, QoS: 1})
		router.RegisterHandler(s.Topic, s.Callback)
	}
	statusTopic := prefs.GetTopicPrefix() + "/status"
	subOpts = append(subOpts, paho.SubscribeOptions{Topic: statusTopic, QoS: 1})
	router.RegisterHandler(statusTopic, func(p *paho.Publish) {
		switch msg := string(p.Payload); msg {
		case "online":
			log.Debug().Msg("Home Assistant online.")
		case "offline":
			log.Debug().Msg("Home Assistant offline.")
		}
	})

	// We will connect to the Eclipse test server (note that you may see messages that other users publish)
	u, err := url.Parse(prefs.GetMQTTServer())
	if err != nil {
		panic(err)
	}

	connOpts := genConnOpts(ctx, u, subOpts, router)

	c, err := autopaho.NewConnection(ctx, connOpts) // starts process; will reconnect until context cancelled
	if err != nil {
		return nil, err
	}
	// Wait for the connection to come up
	if err := c.AwaitConnection(ctx); err != nil {
		return nil, err
	}
	return &Client{conn: c}, nil
}

func genConnOpts(ctx context.Context, server *url.URL, subOpts []paho.SubscribeOptions, router *paho.StandardRouter) autopaho.ClientConfig {
	clientID := "go_hass_anything_" + strconv.Itoa(time.Now().Second())

	connOpts := autopaho.ClientConfig{
		ServerUrls: []*url.URL{server},
		KeepAlive:  20, // Keepalive message should be sent every 20 seconds
		// CleanStartOnInitialConnection defaults to false. Setting this to true will clear the session on the first connection.
		CleanStartOnInitialConnection: false,
		// SessionExpiryInterval - Seconds that a session will survive after disconnection.
		// It is important to set this because otherwise, any queued messages will be lost if the connection drops and
		// the server will not queue messages while it is down. The specific setting will depend upon your needs
		// (60 = 1 minute, 3600 = 1 hour, 86400 = one day, 0xFFFFFFFE = 136 years, 0xFFFFFFFF = don't expire)
		SessionExpiryInterval: 60,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			log.Debug().Msg("MQTT connection up.")
			// Subscribing in the OnConnectionUp callback is recommended (ensures the subscription is reestablished if
			// the connection drops)
			if _, err := cm.Subscribe(ctx, &paho.Subscribe{Subscriptions: subOpts}); err != nil {
				log.Warn().Err(err).Msg("Failed to add subscriptions.")
			}
			log.Debug().Msg("Subscriptions added.")
		},
		OnConnectError: func(err error) { log.Error().Err(err).Msg("Error establishing MQTT connection.") },
		// eclipse/paho.golang/paho provides base mqtt functionality, the below config will be passed in for each connection
		ClientConfig: paho.ClientConfig{
			// If you are using QOS 1/2, then it's important to specify a client id (which must be unique)
			ClientID: clientID,
			// OnPublishReceived is a slice of functions that will be called when a message is received.
			// You can write the function(s) yourself or use the supplied Router
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					log.Trace().Str("topic", pr.Packet.Topic).Msg("Routing message to handler.")
					router.Route(pr.Packet.Packet())
					return true, nil // we assume that the router handles all messages (todo: amend router API)
				},
			},
			OnClientError: func(err error) { log.Error().Err(err).Msg("Client error.") },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Debug().Str("reason", d.Properties.ReasonString).Msg("Server requested disconnect.")
				} else {
					log.Debug().Interface("code", d.ReasonCode).Msg("Server requested disconnect.")
				}
			},
		},
	}
	return connOpts
}

func publish(c *autopaho.ConnectionManager, msgs ...*Msg) error {
	var errs error
	for _, msg := range msgs {
		log.Trace().Str("topic", msg.Topic).Bool("retain", msg.Retained).RawJSON("payload", msg.Message).Msg("Publishing message.")
		// Publish a test message (use PublishViaQueue if you don't want to wait for a response)
		if _, err := c.Publish(context.TODO(), &paho.Publish{
			QoS:     1,
			Topic:   msg.Topic,
			Payload: []byte(msg.Message),
		}); err != nil {
			log.Error().Err(err).Str("topic", msg.Topic).Msg("Error publishing message.")
		}
	}
	return errs
}
