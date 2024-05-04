<!--
 Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
 
 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Go Hass Anything

![MIT](https://img.shields.io/github/license/joshuar/go-hass-anything)
![GitHub last commit](https://img.shields.io/github/last-commit/joshuar/go-hass-anything)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshuar/go-hass-anything?style=flat-square)](https://goreportcard.com/report/github.com/joshuar/go-hass-anything)
[![Go Reference](https://pkg.go.dev/badge/github.com/joshuar/go-hass-anything.svg)](https://pkg.go.dev/github.com/joshuar/go-hass-anything)
[![Release](https://img.shields.io/github/release/joshuar/go-hass-anything?style=flat-square)](https://github.com/joshuar/go-hass-anything/releases/latest)

**Send anything to Home Assistant, through MQTT, powered by Go.**

Go Hass Anything is a framework for writing self-contained apps in Go that can
send data and listen for controls to/from Home Assistant, over MQTT.  This can
be useful for adding sensors or controls to Home Assistant that are not
available through an existing Home Assistant integration.

The code is flexible to be imported as a package into your own Go code to
provide this functionality, or it can be run as its own “agent” process that
will manage any apps you write. See the [example
app](pkg/apps/exampleapp/exampleapp.go) or the [docs](docs/README.md) for more
details on both use-cases.

## ⬇️ Installation

>[!NOTE] It's recommended to create your own container to run a Go Hass Anything
>agent with your own apps. The pre-built container cannot be customised and the
>agent will only run an example app to demonstrate functionality. Create your
>own apps and add them to the agent by following the [development
>docs](docs/development/Apps.md).

## Container

Pre-built containers that can run a demo app can be found on the
[packages](https://github.com/joshuar/go-hass-anything/pkgs/container/go-hass-anything)
page on GitHub.


## 🖱️ Usage

`podman` is the recommended container engine for running Go Hass Anything.

### Configuration

To run the agent, you first need to configure the MQTT connection. Use the command:

```shell
podman run --interactive --tty --rm \
    --volume ~/go-hass-anything:/home/gouser:U \
    ghcr.io/joshuar/go-hass-anything configure
```

This will open a user interface in the terminal to enter MQTT connection
details. You can navigate the fields via the keyboard.

### Running

Once the agent is configured, you can run it. Use the command:

```shell
podman run --name my-go-hass-anything \
    --volume ~/go-hass-anything:/home/gouser:U \
    ghcr.io/joshuar/go-hass-anything
```

This will start the agent and any configured apps. If needed each app will send
initial configuration data to Home Assistant that is needed to register its
entities. After that, the app should start sending
data to Home Assistant and should be visible under the [MQTT
Integration](https://www.home-assistant.io/integrations/mqtt/).

### Other Actions

#### Reset/Remove app data from Home Assistant

If needed/desired, you can remove the app entities from Home Assistant by
running the command:

```shell
podman exec my-go-hass-anything clear
```

After this, there should be no devices (from Go Hass Anything) and associated
entities in Home Assistant. If you want to re-add them, execute the run
command again.

## 🤖 Supported Home Assistant Entities

- [Binary Sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/).
- [Sensor](https://www.home-assistant.io/integrations/sensor.mqtt/).
- [Button](https://www.home-assistant.io/integrations/button.mqtt/).
- [Number](https://www.home-assistant.io/integrations/number.mqtt/).
- [Switch](https://www.home-assistant.io/integrations/switch.mqtt/).

_More to come!_

## 🤝 Compatibility

Currently, only Linux is supported. Though the code is designed to be extensible
to other operating systems. See [the development docs](docs/development/README.md) for
details on how to extend for other operating systems.

## Issues, Feature Requests, Contributing

- Found an issue? Please [report
  it](https://github.com/joshuar/go-hass-anything/issues/new?assignees=joshuar&labels=&template=bug_report.md&title=%5BBUG%5D)!
- Have a suggestion for a feature? Submit a [feature
  request](https://github.com/joshuar/go-hass-anything/issues/new?assignees=joshuar&labels=&template=feature_request.md&title=)!
- Want to help develop Go Hass Anything? See the [contributing
  guidelines](CONTRIBUTING.md).

## License

[MIT](LICENSE)
