// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

import "encoding/json"

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
