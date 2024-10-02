// Copyright (c) 2024 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package cameraapp demonstrates an app that can send images from a connected
// webcam to Home Assistant. It exposes the camera entity (through which images
// are published) as well as start and stop button entities for
// starting/stopping the camera respectively.
//
// It uses the excellent github.com/vladimirvivien/go4vl package for camera
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
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"

	mqtthass "github.com/joshuar/go-hass-anything/v11/pkg/hass"
	mqttapi "github.com/joshuar/go-hass-anything/v11/pkg/mqtt"
)

const (
	appName = "Go Hass Anything Camera Example App"
	appID   = "mqtt_camera_example"
)

// Some defaults for the device file, formats and image size.
var (
	deviceFile    = "/dev/video0"
	preferredFmts = []v4l2.FourCCType{v4l2.PixelFmtMPEG, v4l2.PixelFmtMJPEG, v4l2.PixelFmtJPEG, v4l2.PixelFmtYUYV}
	defaultHeight = 640
	defaultWidth  = 480
)

var ErrUnsupportedOS = errors.New("example camera app only runs on Linux")

// CameraApp is our struct that represents an app to Go Hass Anything.
type CameraApp struct {
	camera      *camera
	images      *mqtthass.ImageEntity
	startButton *mqtthass.ButtonEntity
	stopButton  *mqtthass.ButtonEntity
	msgCh       chan *mqttapi.Msg
}

// camera is an internal struct to hold some data about the camera in use.
type camera struct {
	device     *device.Device
	cancelFunc context.CancelFunc
	fps        time.Duration
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

	app.images = mqtthass.NewImageEntity().
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
				// We are using a function enclosure here for the callback as we
				// need access to the context passed to our app.
				camera, err := openCamera(deviceFile)
				if err != nil {
					slog.Error("Could not open camera device.", slog.Any("error", err))

					return
				}

				app.camera = camera

				camCtx, cancelFunc := context.WithCancel(ctx)
				app.camera.cancelFunc = cancelFunc

				slog.Info("Start recording webcam.")

				go app.camera.publishImages(camCtx, app.images.GetImageTopic(), app.msgCh)
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
				if app.camera.cancelFunc != nil {
					app.camera.cancelFunc()
					slog.Info("Stop recording webcam.")

					if err := app.camera.closeCamera(); err != nil {
						slog.Error("Close camera failed.", slog.Any("error", err))
					}
				}
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
		slog.Error("Could not marshal camera entity config.", "error", err)
	} else {
		configs = append(configs, cameraCfg)
	}

	startButtonCfg, err := a.startButton.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal start button entity config.", "error", err)
	} else {
		configs = append(configs, startButtonCfg)
	}

	stopButtonCfg, err := a.stopButton.MarshalConfig()
	if err != nil {
		slog.Error("Could not marshal start button entity config.", "error", err)
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
		slog.Warn("Unable to marshal start button subscription.", "error", err.Error())
	} else {
		subs = append(subs, startButtonSub)
	}

	stopButtonSub, err := a.stopButton.MarshalSubscription()
	if err != nil {
		slog.Warn("Unable to marshal stop button subscription.", "error", err.Error())
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
func openCamera(cameraDevice string) (*camera, error) {
	camDev, err := device.Open(cameraDevice)
	if err != nil {
		return nil, fmt.Errorf("could not open camera %s: %w", cameraDevice, err)
	}

	fps, err := camDev.GetFrameRate()
	if err != nil {
		return nil, fmt.Errorf("could not determine camera frame rate: %w", err)
	}

	fmtDescs, err := camDev.GetFormatDescriptions()
	if err != nil {
		return nil, fmt.Errorf("could not determine camera formats: %w", err)
	}

	var fmtDesc *v4l2.FormatDescription
	for _, preferredFmt := range preferredFmts {
		fmtDesc = getFormats(fmtDescs, preferredFmt)
		if fmtDesc != nil {
			break
		}
	}

	if fmtDesc == nil {
		return nil, fmt.Errorf("camera does not support any preferred formats: %w", err)
	}

	if err = camDev.SetPixFormat(v4l2.PixFormat{
		Width:       uint32(defaultWidth),
		Height:      uint32(defaultHeight),
		PixelFormat: fmtDesc.PixelFormat,
		Field:       v4l2.FieldNone,
	}); err != nil {
		return nil, fmt.Errorf("could not configure camera: %w", err)
	}

	pixFmt, err := camDev.GetPixFormat()
	if err == nil {
		slog.Info("Camera configured.", slog.Any("format", pixFmt), slog.Any("fps", fps))
	}

	return &camera{
			device: camDev,
			fps:    time.Second / time.Duration(fps),
		},
		nil
}

// publishImages loops over the received frames from the camera and wraps them
// as a MQTT message to be sent back on the bus.
func (c *camera) publishImages(ctx context.Context, topic string, msgCh chan *mqttapi.Msg) {
	if err := c.device.Start(ctx); err != nil {
		slog.Error("Could not start recording", slog.Any("error", err))

		return
	}

	for frame := range c.device.GetOutput() {
		slog.Info("Processing camera frame...")
		// buf := new(bytes.Buffer)
		// encoder := base64.NewEncoder(base64.StdEncoding, buf)
		// encoder.Write(frame)
		// str := base64.StdEncoding.EncodeToString(frame)
		msgCh <- mqttapi.NewMsg(topic, frame)

		time.Sleep(c.fps)
	}
}

// closeCamera wraps the v4l2 camera close method.
func (c *camera) closeCamera() error {
	if err := c.device.Close(); err != nil {
		return fmt.Errorf("could not close camera device: %w", err)
	}

	return nil
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

func getFormats(fmts []v4l2.FormatDescription, pixEncoding v4l2.FourCCType) *v4l2.FormatDescription {
	for _, desc := range fmts {
		if desc.PixelFormat == pixEncoding {
			return &desc
		}
	}

	return nil
}
