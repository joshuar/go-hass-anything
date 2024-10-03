<div align="center">

  <h1>Go Hass Anything</h1>

  <p>
    Send anything to Home Assistant, through MQTT, powered by Go.
  </p>

<p>
<a href="https://github.com/joshuar/go-hass-anything/graphs/contributors">
 <img alt="Contributors Shield" src="https://img.shields.io/github/contributors/joshuar/go-hass-anything">
</a>
<a href="https://github.com/joshuar/go-hass-anything/commits/main/">
 <img alt="Last Update Shield" src="https://img.shields.io/github/last-commit/joshuar/go-hass-anything">
</a>
<a href="https://github.com/joshuar/go-hass-anything/forks">
 <img alt="Forks Shield" src="https://img.shields.io/github/forks/joshuar/go-hass-anything">
</a>
<a href="https://github.com/joshuar/go-hass-anything/stargazers">
 <img alt="Stars Shield" src="https://img.shields.io/github/stars/joshuar/go-hass-anything">
</a>
<a href="https://github.com/joshuar/go-hass-anything/issues">
 <img alt="Open Issues Shield" src="https://img.shields.io/github/issues/joshuar/go-hass-anything">
</a>
<a href="https://github.com/joshuar/go-hass-anything/blob/main/LICENSE">
 <img alt="License Shield" src="https://img.shields.io/github/license/joshuar/go-hass-anything">
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
    - [📦 Cross Compilation](#-cross-compilation)
    - [🚩 Deployment](#-deployment)
  - [🏃 Running](#-running)
    - [🔧 Configuration](#-configuration)
    - [👀 Usage](#-usage)
    - [♻️ Reset](#️-reset)
  - [💻 Development](#-development)
    - [💽 Building Apps](#-building-apps)
      - [Examples](#examples)
      - [Code Location](#code-location)
    - [App Requirements](#app-requirements)
    - [Poll based Apps](#poll-based-apps)
    - [Event based Apps](#event-based-apps)
    - [(Optional) App Configuration](#optional-app-configuration)
    - [Adding to the agent](#adding-to-the-agent)
    - [Logging](#logging)
  - [👋 Contributing](#-contributing)
    - [🏁 Committing Code](#-committing-code)
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

The agent is extremely light on resources, using only a few tens of megabytes of
memory at most. As the agent and apps are written in Go, it can run on a wide
variety of platforms from embedded through to server hardware.

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

### 🎯 Features

- Write self-contained “apps” in Go that are run by the agent.
- Apps can specify either a polling interval that the agent will run the app on
  to publish updates to MQTT, or, pass a channel back to the agent and send
  events that the agent will publish on MQTT.
- Apps can optionally specify user-facing preferences that the agent will
  present via a terminal UI for configuration.
- Apps can use the following Home Assistant entities:
  - [Sensor](https://www.home-assistant.io/integrations/sensor.mqtt/) ([Example App](examples/sensorapp/main.go))
  - [Binary Sensor](https://www.home-assistant.io/integrations/binary_sensor.mqtt/)
  - [Switch](https://www.home-assistant.io/integrations/switch.mqtt/) ([Example App](examples/switchapp/main.go))
  - [Button](https://www.home-assistant.io/integrations/button.mqtt/) ([Example App](examples/buttonapp/main.go))
  - [Number](https://www.home-assistant.io/integrations/number.mqtt/) ([Example App](examples/numberapp/main.go))
  - [Text](https://www.home-assistant.io/integrations/text.mqtt/) ([Example App](examples/textapp/main.go))
  - [Image](https://www.home-assistant.io/integrations/image.mqtt/) ([Example App](examples/cameraapp/main.go))
  - [Camera](https://www.home-assistant.io/integrations/camera.mqtt/)
  - _With more to come!_
- Simple TOML based configuration.
- Compile all apps into a single binary.
- Use via a container or stand-alone binary.
- Light on resources (CPU/memory).
- Runs anywhere that Go runs, from embedded to server hardware.

> [!NOTE]
> If you are looking to add sensors for a Linux desktop/server to Home
> Assistant, check out [Go Hass Agent](https://github.com/joshuar/go-hass-agent)
> which uses Go Hass Anything under the hood!

### 🗒️ Versioning

This project follows [semantic versioning](https://semver.org/). Given a version
number `MAJOR`.`MINOR`.`PATCH`, the gist of it is:

- A `MAJOR` number change means breaking changes from the previous release.
- A `MINOR` number change means significant changes and new features have been
  added, but not breaking changes.
- A `PATCH` number change indicate minor changes and bug fixes.

[⬆️ Back to Top](#-table-of-contents)

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

An example configuration for Mosquitto has been provided in
`deployments/mosquitto/config/mosquitto.conf.example`.

The Mosquitto command-line utilities (`mosquitto_{pub,sub}`) are installed in
the devcontainer.

### ⚙️ Building

> [!NOTE]
> If you have not yet created an app, Go Hass Anything will build with
> an included example app. See the [app creation instructions](#-building-apps)
> below for details on creating and including your own apps.

Use the following mage invocation in the project root directory:

```shell
mage -d build/magefiles -w . build:full
```

This will:

- Run `go generate ./...`.
- Run `go mod tidy`.
- Run `go fmt ./...`.
- Build a binary and place it in `dist/go-hass-anything`.

To just build a binary, replace `build:full` with `build:fast` in the mage
invocation above.

To see all possible build commands, run:

```shell
mage -d build/magefiles -w . -l
```

### 📦 Cross Compilation

Go Hass Anything can also be built for **arm (v6/v7)** and **arm64** with
cross-compilation. To build for a different architecture, set the
`TARGETPLATFORM` environment variable:

```shell
export TARGETPLATFORM=linux/arm64 # or linux/arm/v6 or linux/arm/v7
```

### 🚩 Deployment

While Go Hass Anything can be run as a single binary, using a container is
recommended. `podman` is the container engine of choice for deployment.

A [Dockerfile](./Dockerfile) is available that you can use to build an image
containing your own custom apps.

To add your own apps to the container, copy them into a directory in the base of
the repo (for example, `apps/`) and then specify the build arg `APPDIR` pointing
to this location:

```shell
podman build --file ./Dockerfile --tag go-hass-anything --build-arg APPDIR=apps
```

As with building a binary,
[cross-compliation](https://docs.docker.com/build/building/multi-platform/#cross-compilation)
is supported:

```shell
# use either linux/arm64, linux/arm/v7 or linux/arm/v6
podman build --file ./Dockerfile --platform linux/arm/v7 --tag go-hass-anything --build-arg APPDIR=apps
```

By default, the container will run as a user with uid/gid 1000/1000. You can
pick a different uid/gid when building by adding `--build-arg UID=999` and
`--build-arg GID=999` (adjusting the values as appropriate).

Pre-built containers that can run some demo apps showing some of the available
entities can be found on the
[packages](https://github.com/joshuar/go-hass-anything/pkgs/container/go-hass-anything)
page on GitHub. The demo app source code can be found in
[examples/](./examples/).

[⬆️ Back to Top](#-table-of-contents)

## 🏃 Running

### 🔧 Configuration

To run the agent, you first need to configure the MQTT connection. Use the
command:

```shell
# For containers:
podman run --interactive --tty --rm \
    --volume go-hass-anything:/home/go-hass-anything:U \
    ghcr.io/joshuar/go-hass-anything configure
# For binaries:
go-hass-anything configure
```

This will open a user interface in the terminal to enter MQTT connection details
for the agent, and then any preferences for apps. You can navigate the fields
via the keyboard.

### 👀 Usage

Once the agent is configured, you can run it. Use the command:

```shell
# For containers:
podman run --name my-go-hass-anything \
    --volume go-hass-anything:/home/go-hass-anything:U \
    ghcr.io/joshuar/go-hass-anything
# For binaries:
go-hass-anything run
```

### ♻️ Reset

If needed/desired, you can remove the app entities from Home Assistant by
running the command:

```shell
# For containers:
podman run --interactive --tty --rm \
    --volume ~/go-hass-anything:/home/go-hass-anything:U \
    ghcr.io/joshuar/go-hass-anything clear
# For binaries:
go-hass-anything clear
```

After this, there should be no devices (from Go Hass Anything) and associated
entities in Home Assistant. If you want to re-add them, execute the run
command again.

[⬆️ Back to Top](#-table-of-contents)

## 💻 Development

### 💽 Building Apps

#### Examples

Check out the [examples](./examples/) which a few of the different types of
entities you can create in Home Assistant.

#### Code Location

> [!IMPORTANT]
> The app directory is not committed to version control. This allows your apps to
> remain private. But it also means that if you desire version control of your
> apps, you should set up your own repo for them.

You can put your code in `apps/`. You can create multiple
directories for each app you develop.

> [!NOTE] The filename is important. The generator expects a file named
> `main.go` in app directories for those directories to be considered as an app.
> Make sure you at least have this file if you split your app code into multiple
> files.

### App Requirements

To develop an app to be run by the agent, create a concrete type that satisfies
the `agent.App` interface:

```go
// App represents an app that the agent can run. All apps have the following
// methods, which define how the app should be configured, current states of its
// entities and any subscriptions it wants to watch.
type App interface {
  // Name() is an identifier for the app, used for logging in the agent.
  Name() string
  // Configuration() returns the messages needed to tell Home Assistant how to
  // configure the app and its entities.
  Configuration() []*mqtt.Msg
  // States() are the messages that reflect the app's current state of the
  // entities of the app.
  States() []*mqtt.Msg
  // Subscriptions() are the topics on which the app wants to subscribe and
  // execute a callback in response to a message on that topic.
  Subscriptions() []*mqtt.Subscription
  // Update() is a function that is run at least once by the agent and will
  // usually contain the logic to update the states of all the apps entities.
  // It may be run multiple times, if the app is also considered a polling
  // app. See the definition for PollingApp for details.
  Update(ctx context.Context) error
}
```

- You don't need to worry about setting up a connection to MQTT, the agent will
do that for you.
- `Name()`: This should return the app name as a string. This is used for
defining the app configuration file (if used) and in various places for display
by the agent.
- `Configuration() []*mqtt.Msg`: This function should return an array of
`mqtt.Msg`, each message representing the configuration topics and details for
the sensors provided by the app.
- `States() []*mqtt.Msg`: This function should return an array of `mqtt.Msg`,
each message representing a single state topic for a sensor provided by the app.
- `Subscriptions []*mqtt.Subscription`: This function should return an array of
`mqtt.Subscription`, each message representing a single subscription topic for
which the app wants to listen on. Each of these subscriptions should have a
callback function that is run when a message is received on the topic.
- `Update(ctx context.Context) error`: This function will be called by the agent
at least once. It can be used to update any app state before the agent publishes
app state messages to MQTT. It should respect context cancellation and act
appropriately on this signal.

Create an exported function called `New` that is used to instantiate your app
with the signature:

```go
func New(ctx context.Context) (*yourAppStruct, error)
```

This function should return your concrete type that satisfies the interface
methods above, or an error if the app cannot be initialised. You can put
whatever code you need in this function to set up your application (i.e.,
reading from configs, setting up other connections, etc.). This will be called
first by the agent to initialise your app.

### Poll based Apps

If the app should be run on some kind of interval, updating its state each time,
it should have the following method:

```go
// PollingApp represents an app that should be polled for updates on some
// interval. When an app satisfies this interface, the agent will configure a
// goroutine to run the apps Update() function and publish its States().
type PollingApp interface {
  // PollConfig defines the interval on which the app should be polled and its
  // states updated. A jitter should be defined, that is much less than the
  // interval, to add a small variation to the interval to avoid any
  // "thundering herd" problems.
  PollConfig() (interval, jitter time.Duration)
}
```

### Event based Apps

If the app has its own event loop, and requires states to be published when
certain events occur, it should have the following method:

```go
// EventsApp represents an app that will update its States() in response to some
// event(s) it is monitoring. When an app satisfies this interface, the agent
// will configure a goroutine to watch a channel of messages the app sends when
// an event occurs, which will be published to MQTT.
type EventsApp interface {
  // MsgCh is a channel of messages that the app generates when some internal
  // event occurs and a new message should be published to MQTT.
  MsgCh() chan *mqtt.Msg
}
```

In the app code (usually within `New()`), the app should create a `chan
*mqtt.Msg`, returned by the method above. Any time a state update needs to be
published, it can be sent through this channel and the agent will publish the
message on MQTT.

### (Optional) App Configuration

If your app has user-facing configuration, the agent supports presenting these
to the user when its configuration command is run. It will then create and
utilise a per-app configuration stored in the users home directory
(`~/.config/go-hass-anything/APPNAME-preferences.toml` on Linux).

For your app to support this, make sure it satisfies the `AppWithPreferences`
interface:

```go
// AppWithPreferences represents an app that has preferences that can be
// configured by the user.
type AppWithPreferences interface {
  App
  // DefaultPrefernces returns the AppPreferences map of default preferences for the app.
  // This is passed to the UI code to facilitate generating a form to enter
  // the preferences when the agent runs its configure command.
  DefaultPreferences() (preferences.AppPreferences)
}
```

Each app preference can be represented as a `preference.Preference`:

```go
// Preference represents a single preference in a preferences file.
type Preference struct {
  // Value is the actual preference value.
  Value any `toml:"value"`
  // Description is a string that describes the preference, and may be used
  // for display purposes.
  Description string `toml:"description,omitempty"`
  // Secret is a flag that indicates whether this preference represents a
  // secret. The value has no effect on the preference encoding in the TOML,
  // only on how to display the preference to the user (masked or plaintext).
  Secret bool `toml:"-"`
}
```

The agent takes care of loading and saving the configuration. When the agent is
configured or run for the first time, the agent will show/use default
preferences for each app.

### Adding to the agent

If you have followed the requirements above for both location and code
functions, you can run `go generate ./...` in the repo root to add your app(s)
to the agent. A new `internal/agent/init.go` file should be generated, which
will contain the necessary code to run your apps to the agent.

> [!IMPORTANT]
> The file `internal/agent/init.go` is not committed to version control. Like
> your app code, this allows your apps to remain private.

After building the agent, it should run all of your apps.

### Logging

All packages use [log/slog](https://pkg.go.dev/log/slog) for logging, so if
including the Go Hass Anything packages in your own code, you can hook into
and/or extend upon that. Note that some of the packages define custom levels for
_trace_ (level -8) and _fatal_ (level 12), which if the logger is set to output,
will show some additional details from the internals.

[⬆️ Back to Top](#-table-of-contents)

## 👋 Contributing

<picture>
 <img alt="Go Hass Anything Contributors" src="https://contrib.rocks/image?repo=joshuar/go-hass-anything">
</picture>

Contributions are always welcome!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for ways to get started.

### 🏁 Committing Code

This repository is using [conventional commit
messages](https://www.conventionalcommits.org/en/v1.0.0/#summary). This provides
the ability to automatically include relevant notes in the
[changelog](./CHANGELOG.md). The [TL;DR](https://en.wikipedia.org/wiki/TL;DR) is
when writing commit messages, add a prefix:

- `feat:` for a new feature, like a new sensor.
- `fix:` when fixing an issue.
- `refactor:` when making non-visible but useful code changes.
- …and so on. See the link above or see the existing commit messages for examples.

### 📜 Code of Conduct

Please read the [Code of Conduct](./CODE_OF_CONDUCT.md)

[⬆️ Back to Top](#-table-of-contents)

## ⚠️ License

Distributed under the [MIT](LICENSE) license.

## 🤝 Contact

Joshua Rich - [@joshuar](https://github.com/joshuar)

Project Link: [https://github.com/joshuar/go-hass-anything](https://github.com/joshuar/go-hass-anything)

## 💎 Acknowledgements

- [Shields.io](https://shields.io/)
- [Awesome README](https://github.com/matiassingers/awesome-readme)
- [Emoji Cheat Sheet](https://github.com/ikatyang/emoji-cheat-sheet/blob/master/README.md#travel--places)
- [Home Assistant](https://home-assistant.io).

[⬆️ Back to Top](#-table-of-contents)
