// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mqtt

// Msg represents a message that can be sent or received on the MQTT bus.
type Msg struct {
	Topic    string
	Message  []byte
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
// message body. The returned Msg can be further customized directly for
// specifying retention and QoS parameters, which are not set through this
// function and assumed to be left as their default values.
func NewMsg(topic string, msg []byte) *Msg {
	return &Msg{
		Topic:   topic,
		Message: msg,
	}
}
