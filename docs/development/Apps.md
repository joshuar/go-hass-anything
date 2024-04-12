<!--
 Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>

 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Developing Apps

You can develop an app to send data to Home Assistant via MQTT in two ways:

- As a stand-alone program, importing and using the modules under `pkg/`.
- As a group of apps managed by the *agent* by satisfying an interface defined
  by the agent needed to manage each app.

Which path to take is up to you.  If you want to send different data from
different places to Home Assistant, you'll likely want to develop apps that run
under the agent as a single executable. Otherwise, you can embed the packages
under `pkg/` into your own application.

In either case, apps should satisfy [Home Assistant MQTT
Discovery](https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery).
Their operation is roughly:

1. Send discovery message/payload (once only).
2. Send state messages/payloads (as required).

For an app using Go Hass Anything, that means the app satisfies the
`mqtt.Device` and `Agent.App` interfaces. These share common methods, so it is
easy to satisfy both. 

If the app is self-contained and not running through the agent framework, it
should satisfy `mqtt.Device`.

## Example App running under the Agent

Check out the [example app](../../internal/apps/exampleApp/exampleApp.go) which
demonstrates:

- How to get data from the system running Go Hass anything using an external Go
  package.
- How to get data from the web using a helper function to issue web requests.
- How to create a button in Home Assistant that when pressed, will execute an
  app action.

The code has lots of comments for guidance. It can also be heavily optimised
from its current state and so is a good starting point for practising Go as
well.

## Developing for the agent

### Code Location

> [!IMPORTANT]
> The app directory is not committed to version control. This allows your apps to
> remain private. But it also means that if you desire version control of your
> apps, you should set up your own repo for them.

You can put your code in `apps`. You can create multiple
directories for each app you develop.

> [!NOTE]
> The filename is important. The generator to automatically add your app
> to the agent needs a `.go` file named the same as the app directory to detect
> your app. Make sure you at least have this file if you split your app code
> into multiple files.

### Configuration

When using the agent, it will create and utilise per-app configurations stored
in the users home directory
(`~/.config/go-hass-anything/APPNAME-preferences.toml` on Linux). Within your
app code, you can retrieve the config with
`preferences.LoadAppPreferences(appname)`, which will return a
`preferences.AppConfig` containing a `map[string]any` called `Prefs`, containing
any stored preferences as a key-value map.

### Code Requirements

To develop an app to be run by the agent, create a concrete type that satisfies
the `mqtt.Device` and `agent.App`. Effectively, the type should have the
following methods:

```go
	Name() string
	Configuration() []*mqtt.Msg
	States() []*mqtt.Msg
	Subscriptions() []*mqtt.Subscription
	Run(ctx context.Context, client hass.MQTTClient) error
```

You don't need to worry about setting up a connection to MQTT (satisfying
`mqtt.Client`), the agent will do that for you.

The `Run` function should respect context cancellation and act appropriately on
this signal.

Create a function called `New` that is used to instantiate your app with the signature:

```go
func New(ctx context.Context) (*yourAppStruct, error)
```

This function should return your concrete type that satisfies the interface
methods above, or an error if the app cannot be initialised. You can put
whatever code you need in this function to set up your application (i.e.,
reading from configs, setting up connections, etc.).

The interface methods are explained below.

#### Name()

This should return the app name as a string. This is used for defining the
app configuration file (if used) and in various places for display by the agent.


#### Configuration() []*mqtt.Msg

This function should return an array of `mqtt.Msg`, each message representing
the configuration topics and details for the sensors provided by the app.

#### States() []*mqtt.Msg

This function should return an array of `mqtt.Msg`, each message representing a
single state topic for a sensor provided by the app.

#### Subscriptions []*mqtt.Subscription

This function should return an array of `mqtt.Subscription`, each message representing a
single subscription topic for which the app wants to listen on. Each of these
subscriptions should have a callback function that is run when a message is
received on the topic.

#### Run(ctx context.Context, client *mqtt.Client) error

This function should, on some kind of interval or event/feedback loop, run
`client.Publish(msgs ...*mqtt.Msg)` to send state payloads to Home Assistant. It
can also set up any subscriptions directly. If you concrete type for the app
satisfies the `mqtt.Device` interface, then this call would be
`client.Publish(yourType.States()...)`.

### Adding to the agent

If you have followed the requirements above for both location and code
functions, you can run `go generate ./...` in the repo root to add your app(s)
to the agent. A new `internal/agent/init.go` file should be generated, which
will contain the necessary code to add your apps to the agent.

> [!IMPORTANT]
> The file `internal/agent/init.go` is not committed to version control. Like
> your app code, this allows your apps to remain private.

After running `go build` (or building a release with `goreleaser`) the agent
should run all of your apps.
