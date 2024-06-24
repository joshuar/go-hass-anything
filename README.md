<div align="center">

  <img src="assets/logo.png" alt="logo" width="200" height="auto" />
  <h1>Go Hass Anything</h1>

  <p>
    Send anything to Home Assistant, through MQTT, powered by Go.
  </p>

<p>
  <a href="https://github.com/joshuar/go-hass-anything/graphs/contributors">
    <img src="https://img.shields.io/github/contributors/joshuar/go-hass-anything" alt="contributors" />
  </a>
  <a href="">
    <img src="https://img.shields.io/github/last-commit/joshuar/go-hass-anything" alt="last update" />
  </a>
  <a href="https://github.com/joshuar/go-hass-anything/network/members">
    <img src="https://img.shields.io/github/forks/joshuar/go-hass-anything" alt="forks" />
  </a>
  <a href="https://github.com/joshuar/go-hass-anything/stargazers">
    <img src="https://img.shields.io/github/stars/joshuar/go-hass-anything" alt="stars" />
  </a>
  <a href="https://github.com/joshuar/go-hass-anything/issues/">
    <img src="https://img.shields.io/github/issues/joshuar/go-hass-anything" alt="open issues" />
  </a>
  <a href="https://github.com/joshuar/go-hass-anything/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/joshuar/go-hass-anything.svg" alt="license" />
  </a>
</p>

<h4>
    <a href="https://github.com/joshuar/go-hass-anything">Documentation</a>
  <span> · </span>
    <a href="https://github.com/joshuar/go-hass-anything/issues/">Report Bug</a>
  <span> · </span>
    <a href="https://github.com/joshuar/go-hass-anything/issues/">Request Feature</a>
  </h4>
</div>

<br />

<!-- Table of Contents -->
# 📓 Table of Contents

- [📓 Table of Contents](#-table-of-contents)
  - [🌟 About the Project](#-about-the-project)
    - [👾 Tech Stack](#-tech-stack)
    - [🎯 Features](#-features)
    - [🗒️ Versioning](#️-versioning)
  - [🧰 Getting Started](#-getting-started)
    - [‼️ Prerequisites](#️-prerequisites)
    - [🚧 Development Environment](#-development-environment)
    - [⚙️ Building](#️-building)
    - [🚩 Deployment](#-deployment)
    - [🏃 Running](#-running)
      - [🔧 Configuration](#-configuration)
      - [👀 Usage](#-usage)
      - [♻️ Reset](#️-reset)
  - [👋 Contributing](#-contributing)
    - [📜 Code of Conduct](#-code-of-conduct)
  - [⚠️ License](#️-license)
  - [🤝 Contact](#-contact)
  - [💎 Acknowledgements](#-acknowledgements)

## 🌟 About the Project

Go Hass Anything is a framework for writing self-contained apps in Go that can
send data and listen for controls to/from Home Assistant, over MQTT.  This can
be useful for adding sensors or controls to Home Assistant that are not
available through an existing Home Assistant integration.

The code is flexible to be imported as a package into your own Go code to
provide this functionality, or it can be run as its own “agent” process that
will manage any apps you write.

### 👾 Tech Stack

<details>
  <summary>Agent</summary>
  <ul>
    <li><a href="https://go.dev/">Golang powered</a></li>
    <li><a href="https://github.com/pelletier/go-toml">TOML configuration files</a></li>
    <li><a href="https://github.com/charmbracelet/bubbletea">BubbleTea UI</a></li>
    <li><a href="https://github.com/eclipse/paho.golang">Eclipse Pahoe MQTT Client</a></li>
  </ul>
</details>

<details>
<summary>DevOps</summary>
  <ul>
    <li><a href="https://magefile.org/">Mage for build system.</a></li>
    <li><a href="https://github.com/">GitHub for SCM, CI.</a></li>
  </ul>
</details>

<!-- Features -->
### 🎯 Features

- Write self-contained “apps” in Go that are run by the agent.
- Apps can specify either a polling interval that the agent will run the app on
  to publish updates to MQTT, or, pass a channel back to the agent and send
  events that the agent will publish on MQTT.
- Apps can optionally specify user-facing preferences that the agent will
  present via a terminal UI for configuration.
- Apps can use the following Home Assistant entities:
  - [Sensor](https://www.home-assistant.io/integrations/sensor.mqtt/)
  - [Binary Sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/)
  - [Switch](https://www.home-assistant.io/integrations/switch.mqtt/)
  - [Button](https://www.home-assistant.io/integrations/button.mqtt/)
  - [Number](https://www.home-assistant.io/integrations/number.mqtt/)
  - _With more to come!_
- Simple TOML based configuration.

### 🗒️ Versioning

This project follows [semantic versioning](https://semver.org/). Given a version
number `MAJOR`.`MINOR`.`PATCH`, the gist of it is:

- A `MAJOR` number change means breaking changes from the previous release.
- A `MINOR` number change means significant changes and new features have been
  added, but not breaking changes.
- A `PATCH` number change indicate minor changes and bug fixes.

## 🧰 Getting Started

### ‼️ Prerequisites

Go Hass Anything uses [Mage](https://magefile.org/) for development. Make sure you
follow the instructions on the Mage website to install Mage. If you are using
the devcontainer (see below), this is already installed.

### 🚧 Development Environment

It is recommended to use [Visual Studio Code](https://code.visualstudio.com/).
This project makes use of a [Devcontainer](https://containers.dev/) to provide
some convenience during development.

[![Open in Dev Containers](https://img.shields.io/static/v1?label=Dev%20Containers&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/joshuar/go-hass-anything)

If using Visual Studio Code, you should be prompted when opening your cloned
copy of the code to set up the dev container. The container contains an
installation of Home Assistant and Mosquitto (MQTT broker) that can be used for
testing. They should be started automatically.

- Home Assistant will be listening on <http://localhost:8123>.
- Mosquitto will be listening on <http://localhost:1833>.

### ⚙️ Building

> [!NOTE]
> If you have not yet created an app, Go Hass Anything will build with
> an included example app. See the app creation instructions below for details
> on creating and including your own apps.

Use the following mage invocation in the project root directory:

```shell
mage -d build/magefiles -w . build:full
```

This will:

- Run `go mod tidy`.
- Run `go fmt ./...`.
- Run `go generate ./...`.
- Build a binary and place it in `dist/go-hass-anything`.

To just build a binary, replace `build:full` with `build:fast` in the mage
invocation above.

To see all possible build commands, run:

```shell
mage -d build/magefiles -w . -l
```

### 🚩 Deployment

While Go Hass Anything can be run as a single binary, using a container is recommended. `podman` is the container engine of choice for deployment.

A [Dockerfile](./Dockerfile) is available that you can use to build an image
containing your own custom apps.

To add your own apps to the container, copy them into a directory in the base of
the repo (for example, `apps/`) and then specify the build arg `APPDIR` pointing
to this location:

```shell
podman build --file ./Dockerfile --tag go-hass-anything --build-arg APPDIR=apps
```

Pre-built containers that can run a demo app can be found on the
[packages](https://github.com/joshuar/go-hass-anything/pkgs/container/go-hass-anything)
page on GitHub.

<!-- Run Locally -->
### 🏃 Running

#### 🔧 Configuration

To run the agent, you first need to configure the MQTT connection. Use the
command:

```shell
# For containers:
podman run --interactive --tty --rm \
    --volume ~/go-hass-anything:/home/gouser:U \
    ghcr.io/joshuar/go-hass-anything configure
# For binaries:
go-hass-anything configure
```

This will open a user interface in the terminal to enter MQTT connection
details for the agent, and then any preferences for apps. You can navigate the fields via the keyboard.

#### 👀 Usage

Once the agent is configured, you can run it. Use the command:

```shell
# For containers:
podman run --name my-go-hass-anything \
    --volume ~/go-hass-anything:/home/gouser:U \
    ghcr.io/joshuar/go-hass-anything
# For binaries:
go-hass-anything run
```

#### ♻️ Reset

If needed/desired, you can remove the app entities from Home Assistant by
running the command:

```shell
# For containers:
podman exec my-go-hass-anything clear
# For binaries:
go-hass-anything clear
```

After this, there should be no devices (from Go Hass Anything) and associated
entities in Home Assistant. If you want to re-add them, execute the run
command again.

<!-- Contributing -->
## 👋 Contributing

<a href="https://github.com/joshuar/go-hass-anything/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=joshuar/go-hass-anything" />
</a>

Contributions are always welcome!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for ways to get started.

### 📜 Code of Conduct

Please read the [Code of Conduct](./CODE_OF_CONDUCT.md)

## ⚠️ License

Distributed under the [MIT](LICENSE) license.

## 🤝 Contact

Joshua Rich - [@joshuar](https://github.com/joshuar)

Project Link: [https://github.com/joshuar/go-hass-anything](https://github.com/joshuar/go-hass-anything)

## 💎 Acknowledgements

Use this section to mention useful resources and libraries that you have used in your projects.

- [Shields.io](https://shields.io/)
- [Awesome README](https://github.com/matiassingers/awesome-readme)
- [Emoji Cheat Sheet](https://github.com/ikatyang/emoji-cheat-sheet/blob/master/README.md#travel--places)
- [Home Assistant](https://home-assistant.io).
