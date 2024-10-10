// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package cameraapp demonstrates an app that can send images from a connected
// webcam to Home Assistant. It exposes the camera entity (through which images
// are published) as well as start and stop button entities for
// starting/stopping the camera respectively.
//
// It uses the excellent https://github.com/blackjack/webcam package for camera
// access.
//
// Unfortunately, this example app only runs on Linux.
package cameraapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"slices"

	"github.com/blackjack/webcam"
	"github.com/eclipse/paho.golang/paho"

	mqtthass "github.com/joshuar/go-hass-anything/v11/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

const (
	appName = "Go Hass Anything Camera Example App"
	appID   = "mqtt_camera_example"
)

// Some defaults for the device file, formats and image size.
var (
	deviceFile    = "/dev/video1"
	preferredFmts = []string{"Motion-JPEG"}
	defaultHeight = 640
	defaultWidth  = 480
)

var ErrUnsupportedOS = errors.New("example camera app only runs on Linux")

// CameraApp is our struct that represents an app to Go Hass Anything.
type CameraApp struct {
	camera      *webcam.Webcam
	images      *mqtthass.CameraEntity
	startButton *mqtthass.ButtonEntity
	stopButton  *mqtthass.ButtonEntity
	msgCh       chan *mqttapi.Msg
}

// New sets up our example app. It creates entities for the camera images and
// start and stop buttons.
func New(ctx context.Context) (*CameraApp, error) {
	if runtime.GOOS != "linux" {
		return nil, ErrUnsupportedOS
	}

	app := &CameraApp{
		msgCh: make(chan *mqttapi.Msg),
	}

	app.images = mqtthass.NewCameraEntity().
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Webcam"),
			mqtthass.ID("webcam"),
			mqtthass.DeviceInfo(newDevice()),
		)

	app.startButton = mqtthass.NewButtonEntity().
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Start Webcam"),
			mqtthass.ID("start_webcam"),
			mqtthass.DeviceInfo(newDevice()),
			mqtthass.Icon("mdi:play"),
		).
		WithCommand(
			mqtthass.CommandCallback(func(_ *paho.Publish) {
				// Open the camera device.
				camera, err := openCamera(deviceFile)
				if err != nil {
					slog.Error("Could not open camera device.",
						slog.Any("error", err))
					return
				}

				app.camera = camera

				slog.Info("Start recording webcam.")

				go publishImages(camera, app.images.Topic, app.msgCh)
			}),
		)

	app.stopButton = mqtthass.NewButtonEntity().
		WithDetails(
			mqtthass.App(appName),
			mqtthass.Name("Stop Webcam"),
			mqtthass.ID("stop_webcam"),
			mqtthass.DeviceInfo(newDevice()),
			mqtthass.Icon("mdi:stop"),
		).
		WithCommand(
			mqtthass.CommandCallback(func(_ *paho.Publish) {
				if err := app.camera.StopStreaming(); err != nil {
					slog.Error("Stop streaming failed.", slog.Any("error", err))
				}

				if err := app.camera.Close(); err != nil {
					slog.Error("Close camera failed.", slog.Any("error", err))
				}

				slog.Info("Stop recording webcam.")
			}),
		)

	go func() {
		defer close(app.msgCh)

		<-ctx.Done()
	}()

	return app, nil
}

//revive:disable:unused-receiver
func (a *CameraApp) Name() string {
	return appName
}

//nolint:mnd
func (a *CameraApp) Configuration() []*mqttapi.Msg {
	configs := make([]*mqttapi.Msg, 0, 3)

	cameraCfg, err := a.images.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal camera entity config.",
			slog.Any("error", err))
	} else {
		configs = append(configs, cameraCfg)
	}

	startButtonCfg, err := a.startButton.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal start button entity config.",
			slog.Any("error", err))
	} else {
		configs = append(configs, startButtonCfg)
	}

	stopButtonCfg, err := a.stopButton.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal start button entity config.",
			slog.Any("error", err))
	} else {
		configs = append(configs, stopButtonCfg)
	}

	return configs
}

// States is unused in this example app. This could be used to take an image
// snapshot via the webcam or expose other camera information as entities...
func (a *CameraApp) States() []*mqttapi.Msg { return nil }

//nolint:mnd
func (a *CameraApp) Subscriptions() []*mqttapi.Subscription {
	subs := make([]*mqttapi.Subscription, 0, 2)

	startButtonSub, err := a.startButton.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal start button subscription.",
			slog.Any("error", err))
	} else {
		subs = append(subs, startButtonSub)
	}

	stopButtonSub, err := a.stopButton.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal stop button subscription.",
			slog.Any("error", err))
	} else {
		subs = append(subs, stopButtonSub)
	}

	return subs
}

// Update is unused, there is no app data to update.
func (a *CameraApp) Update(_ context.Context) error { return nil }

// Our channel on which we send camera images as MQTT messages.
func (a *CameraApp) MsgCh() chan *mqttapi.Msg {
	return a.msgCh
}

// openCamera opens the camera device and ensures that it has a preferred image
// format, framerate and dimensions.
func openCamera(cameraDevice string) (*webcam.Webcam, error) {
	cam, err := webcam.Open(cameraDevice)
	if err != nil {
		return nil, fmt.Errorf("could not open camera %s: %w", cameraDevice, err)
	}

	// select pixel format
	var preferredFormat webcam.PixelFormat

	for format, desc := range cam.GetSupportedFormats() {
		if slices.Contains(preferredFmts, desc) {
			preferredFormat = format
			break
		}
	}

	if preferredFormat == 0 {
		return nil, errors.New("could not determine an appropriate format")
	}

	_, _, _, err = cam.SetImageFormat(preferredFormat, uint32(defaultWidth), uint32(defaultHeight))
	if err != nil {
		return nil, fmt.Errorf("could not set camera parameters: %w", err)
	}

	return cam, nil
}

// publishImages loops over the received frames from the camera and wraps them
// as a MQTT message to be sent back on the bus.
func publishImages(cam *webcam.Webcam, topic string, msgCh chan *mqttapi.Msg) {
	if err := cam.StartStreaming(); err != nil {
		slog.Error("Could not start recording", slog.Any("error", err))

		return
	}

	for {
		err := cam.WaitForFrame(uint32(5))
		if err != nil && errors.As(err, &webcam.Timeout{}) {
			continue
		}

		frame, err := cam.ReadFrame()
		if len(frame) == 0 || err != nil {
			break
		}

		msgCh <- mqttapi.NewMsg(topic, frame)
	}
}

func newDevice() *mqtthass.Device {
	return &mqtthass.Device{
		Name:         appName,
		Identifiers:  []string{appID},
		URL:          "https://github.com/joshuar/go-hass-anything",
		Manufacturer: "go-hass-anything",
		Model:        appID,
	}
}
