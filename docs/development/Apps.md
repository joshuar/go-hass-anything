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

Which path to take is up to you. As a suggestion, if you want to send different
data from different places to Home Assistant, you'll likely want to develop apps
that run under the agent as a single executable, as opposed to having multiple
separate Go executables for each app and its data.

In either case, apps should satisfy [Home Assistant MQTT
Discovery](https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery).
Their operation is roughly:

1. Send discovery message/payload (once only).
2. Send state messages/payloads (as required).

For an app using Go Hass Anything, that means the app:

- Satisfies the `hass.MQTTDevice` interface.
- Uses `hass.PublishConfigs` to send its discovery message.
- Uses `hass.PublishState` to send its state messages.

Additionally, if the app is self-contained and not running through the agent
framework, it will need to supply a `hass.MQTTClient` to the `hass.*` functions.

## Example App

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
> App directories are not committed to version control. This allows your apps to
> remain private. But it also means that if you desire version control of your
> apps, you should set up either git submodules or use subtree merging.

You can put your code in `internal/apps/myapp/myapp.go`. You can create multiple
directories for each app you develop.

> [!NOTE]
> The filename is important. The generator to automatically add your app
> to the agent needs a `.go` file named the same as the app directory to detect
> your app. Make sure you at least have this file if you split your app code
> into multiple files.

### Code Requirements

To develop an app to be run by the agent:

- Create a concrete type that satisfies the `hass.MQTTDevice` interface.
- Create the following exported functions:
  - `Run(context.Context, hass.MQTTClient)`.
  - `Clear(context.Context, hass.MQTTClient)`.

You don't need to worry about setting up a connection to MQTT (satisfying
`hass.MQTTClient`), the agent will do that for you.

All functions should respect context cancellation and act appropriately on this
signal.

You need to at least have

#### Run Function

This function should:

- Run `hass.PublishConfigs(hass.MQTTDevice, hass.MQTTClient)` **once** to send
  the config messages and payloads to Home Assistant. Pass your concrete type
  (that satisfies `hass.MQTTDevice`) and the `hass.MQTTClient` passed in by the
  agent.
- On some kind of interval or event/feedback loop, run
  `hass.PublishState(hass.MQTTDevice, hass.MQTTClient)` to send state payloads
  to Home Assistant.

#### Clear Function

This function should:

- Run `hass.UnPublish(hass.MQTTDevice, hass.MQTTClient)` to remove all entities
  and the device associated with this app from Home Assistant. Repeat runs
  should be a no-op.

### Adding to the agent

Apps the agent runs need their `Run` function appended to the agent `RunList`
variable and their `Clear` function appended to the agent `ClearList` variable.

If you have followed the requirements above for both location and code
functions, you can run `go generate ./...` in the repo root to add your app(s)
to the agent. A new `internal/agent/init.go` file should be generated, which
will contain the necessary code to add your apps to the agent.

> [!IMPORTANT]
> The file `internal/agent/init.go` is not committed to version control. Like
> your app code, this allows your apps to remain private.

After running `go build` (or building a release with `goreleaser`) the agent
should run all of your apps.
