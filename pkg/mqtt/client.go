// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

const (
	defaultKeepAliveSec     = 20
	defaultSessionExpirySec = 60

	// clientIDPrefix is prepended to the generated MQTT client ID.
	clientIDPrefix = "go_hass_anything_"
	// clientIDMaxHostLen is the longest hostname that is embedded in a client ID
	// verbatim. Longer hostnames are hashed instead, which keeps the client ID a
	// sensible length without giving up uniqueness.
	clientIDMaxHostLen = 32
)

var (
	ErrNoConnection       = errors.New("no MQTT connection")
	ErrInvalidTopicPrefix = errors.New("invalid topic prefix")
	ErrInvalidServer      = errors.New("invalid server")
	ErrNoPrefs            = errors.New("no preferences provided")
)

type Preferences interface {
	TopicPrefix() string
	Server() string
	User() string
	Password() string
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
	conn     *autopaho.ConnectionManager
	haStatus chan string
}

// Publish will send the list of messages it is passed to the broker that the
// client is connected to. Any errors in publihsing will be returned.
func (c *Client) Publish(ctx context.Context, msgs ...*Msg) error {
	if c.conn == nil {
		return ErrNoConnection
	}

	err := publish(ctx, c.conn, msgs...)

	return err
}

func (c *Client) Unpublish(ctx context.Context, msgs ...*Msg) error {
	if c.conn == nil {
		return ErrNoConnection
	}

	newMsgs := make([]*Msg, 0, len(msgs))

	for _, msg := range msgs {
		newMsgs = append(newMsgs, NewMsg(msg.Topic, []byte(``)))
	}

	err := publish(ctx, c.conn, newMsgs...)

	return err
}

//nolint:exhaustruct
func NewClient(ctx context.Context, prefs Preferences, subscriptions []*Subscription, configs []*Msg) (*Client, error) {
	if prefs == nil {
		return nil, ErrNoPrefs
	}

	subOpts := make([]paho.SubscribeOptions, 0, len(subscriptions))

	client := &Client{
		haStatus: make(chan string),
	}
	router := paho.NewStandardRouter()

	for _, s := range subscriptions {
		slog.Debug("Adding subscription for topic.",
			slog.String("topic", s.Topic))

		subOpts = append(subOpts, paho.SubscribeOptions{Topic: s.Topic, QoS: 1})
		router.RegisterHandler(s.Topic, s.Callback)
	}

	if prefs.TopicPrefix() == "" {
		return nil, ErrInvalidTopicPrefix
	}

	statusTopic := prefs.TopicPrefix() + "/status"
	subOpts = append(subOpts, paho.SubscribeOptions{Topic: statusTopic, QoS: 1})
	router.RegisterHandler(statusTopic, func(p *paho.Publish) {
		client.haStatus <- string(p.Payload)
	})

	connOpts := genConnOpts(ctx, prefs, subOpts, router)

	conn, err := autopaho.NewConnection(ctx, connOpts) // starts process; will reconnect until context canceled.
	if err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}
	// Wait for the connection to come up
	if err := conn.AwaitConnection(ctx); err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}

	client.conn = conn

	if err := client.Publish(ctx, configs...); err != nil {
		slog.Error("Failed to publish configuration messages.",
			slog.Any("error", err))
	}

	client.monitorHAStatus(ctx, configs...)

	return client, nil
}

//nolint:exhaustruct
func genConnOpts(ctx context.Context, prefs Preferences, subOpts []paho.SubscribeOptions, router *paho.StandardRouter) autopaho.ClientConfig { //nolint:lll
	// Set a client ID for this connection.
	clientID := newClientID()

	// Get the server from the preferences and convert to a URL.
	serverURL, err := url.Parse(prefs.Server())
	if err != nil {
		panic(err)
	}

	connOpts := autopaho.ClientConfig{
		ServerUrls: []*url.URL{serverURL},
		KeepAlive:  defaultKeepAliveSec, // Keepalive message should be sent every 20 seconds
		// CleanStartOnInitialConnection defaults to false. Setting this to true will clear the session on the first connection.
		CleanStartOnInitialConnection: false,
		// SessionExpiryInterval - Seconds that a session will survive after disconnection.
		// It is important to set this because otherwise, any queued messages will be lost if the connection drops and
		// the server will not queue messages while it is down. The specific setting will depend upon your needs
		// (60 = 1 minute, 3600 = 1 hour, 86400 = one day, 0xFFFFFFFE = 136 years, 0xFFFFFFFF = don't expire)
		SessionExpiryInterval: defaultSessionExpirySec,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, _ *paho.Connack) {
			slog.Debug("MQTT connection up.")
			// Subscribing in the OnConnectionUp callback is recommended (ensures the subscription is reestablished if
			// the connection drops)
			if _, err := cm.Subscribe(ctx, &paho.Subscribe{Subscriptions: subOpts}); err != nil {
				slog.Warn("Failed to publish subscriptions to MQTT.",
					slog.Any("error", err))
			}
			slog.Debug("Subscriptions added to MQTT.")
		},
		OnConnectError: func(err error) {
			slog.Error("Error establishing MQTT connection.",
				slog.Any("error", err))
		},
		// eclipse/paho.golang/paho provides base mqtt functionality, the below config will be passed in for each connection
		ClientConfig: paho.ClientConfig{
			// If you are using QOS 1/2, then it's important to specify a client id (which must be unique)
			ClientID: clientID,
			// OnPublishReceived is a slice of functions that will be called when a message is received.
			// You can write the function(s) yourself or use the supplied Router
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					slog.Log(ctx, LevelTrace, "Routing message to handler.",
						slog.String("topic", pr.Packet.Topic))
					router.Route(pr.Packet.Packet())

					return true, nil // we assume that the router handles all messages (todo: amend router API)
				},
			},
			OnClientError: func(err error) {
				slog.Error("Client error.",
					slog.Any("error", err))
			},
			OnServerDisconnect: func(packet *paho.Disconnect) {
				if packet.Properties != nil {
					slog.Debug("Server requested disconnect.",
						slog.String("reason", packet.Properties.ReasonString))
				} else {
					slog.Debug("Server requested disconnect.",
						slog.Any("code", packet.ReasonCode))
				}
			},
		},
	}

	// If a username/password is set, add those to the connection options.
	if prefs.User() != "" && prefs.Password() != "" {
		connOpts.ConnectUsername = prefs.User()
		connOpts.ConnectPassword = []byte(prefs.Password())
	}

	return connOpts
}

// newClientID generates the MQTT client ID to connect with. The ID needs to be
// both unique to this client and stable across reconnections:
//
//   - MQTT requires client IDs to be unique. When a client connects with an ID
//     that is already connected, the broker closes the existing connection, so
//     two clients sharing an ID will repeatedly evict one another.
//   - The broker keys the session state it retains on the client ID, and this
//     client asks for a session that outlives disconnection (see
//     SessionExpiryInterval above). An ID that changes between reconnections
//     therefore orphans the session belonging to the previous ID.
//
// Deriving the ID from the hostname satisfies both properties for the usual
// case of one instance per device. Where the hostname cannot be used, a random
// ID is generated instead, which is unique but not stable.
func newClientID() string {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Debug("Could not determine hostname for MQTT client ID, generating a random one.",
			slog.Any("error", err))
	}

	return clientIDForHost(hostname)
}

// clientIDForHost derives the client ID for the given hostname. An unusable
// hostname yields a random ID, and an overly long one is hashed rather than
// truncated, so that distinct hostnames always yield distinct client IDs.
func clientIDForHost(hostname string) string {
	id := sanitiseClientID(hostname)

	switch {
	case id == "":
		id = rand.Text()
	case len(id) > clientIDMaxHostLen:
		sum := sha256.Sum256([]byte(id))
		id = hex.EncodeToString(sum[:8])
	}

	return clientIDPrefix + id
}

// sanitiseClientID drops any characters from the given string that cannot be
// relied upon in an MQTT client ID. MQTT only requires brokers to accept
// alphanumeric client IDs, so anything else is removed.
func sanitiseClientID(id string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '0' && r <= '9':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		default:
			return -1
		}
	}, id)
}

//nolint:exhaustruct
func publish(ctx context.Context, conn *autopaho.ConnectionManager, msgs ...*Msg) error {
	var errs error

	for _, msg := range msgs {
		slog.Log(ctx, LevelTrace, "Publishing message.",
			slog.String("topic", msg.Topic),
			slog.Bool("retain", msg.Retained),
			slog.Any("payload", msg.Message))

		if _, err := conn.Publish(ctx, &paho.Publish{
			QoS:     1,
			Topic:   msg.Topic,
			Payload: msg.Message,
		}); err != nil {
			slog.Error("Error publishing message.",
				slog.String("topic", msg.Topic),
				slog.Any("error", err))
		}
	}

	return errs
}

func (c *Client) monitorHAStatus(ctx context.Context, configs ...*Msg) {
	go func() {
		for {
			select {
			case status := <-c.haStatus:
				switch status {
				case "online":
					slog.Debug("Home Assistant detected online.")

					if err := c.Publish(ctx, configs...); err != nil {
						slog.Warn("Could not publish configs to MQTT",
							slog.Any("error", err))
					}
				case "offline":
					slog.Debug("Home Assistant detected offline.")
				}
			case <-ctx.Done():
				close(c.haStatus)

				return
			}
		}
	}()
}
