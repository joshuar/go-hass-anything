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

## ⬇️ Installation

## Packages

Head over to the [releases](https://github.com/joshuar/go-hass-anything/releases)
page and download the appropriate package for your operating system and/or
distribution:

- For **Fedora**, use the `.rpm`.
- For **Ubuntu**, use the `.deb`.
- For **Arch**, use the `.tar.zst`.

Any other distributions not listed above, you can try the binary, or build it
yourself from source (see [development docs](docs/development/README.md)).

## Container

A Dockerfile that you can use to build an image can be found [here](Dockerfile).

To add your own apps to the container, copy them into a directory in the base of
the repo (for example, `apps/`) and then specify the build arg `APPDIR` pointing
to this location:

```shell
podman build --file ./Dockerfile --tag go-hass-anything --build-arg APPDIR=apps
```

## 🖱️ Usage

### Configuration

To run the agent, you first need to configure the MQTT connection. For a package install, use the command:

```shell
go-hass-anything configure
```

For a container, use the following:

```shell
podman run --interactive --tty --rm \
    --volume ~/go-hass-anything:/home/gouser:U \
    ghcr.io/joshuar/go-hass-anything configure
```

This will open a user interface in the terminal to enter MQTT connection
details. You can navigate the fields via the keyboard.

### Running

Once the agent is configured, you can run it. To do so, for package installs, use the command:

```shell
go-hass-anything run
```

For a container, run:

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

> [!NOTE]
> By default, the agent only runs a single [example
> app](internal/apps/exampleApp/exampleApp.go) that will fetch load averages from
> the local system, current temperature from an online weather service and add a
> button to open the Home Assistant homepage. Add your own apps by following the
> [development docs](docs/development/README.md).

### Other Actions

#### Reset/Remove app data from Home Assistant

If needed/desired, you can remove the app entities from Home Assistant by
running the command:

```shell
# for a package install
go-hass-anything clear
# for a container install
podman exec my-go-hass-anything clear
```

After this, there should be no devices (from Go Hass Anything) and associated
entities in Home Assistant. If you want to re-add them, execute the run
command again.

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
