<!--
 Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
 
 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Go Hass Anything Docs

## Concepts

- **Agent** refers to the `go-hass-anything` executable that will handle connecting to MQTT and managing/running a bunch of *Apps*.
- **App** refers to an *application* that gets data and asks the agent to send it to Home Assistant or sends the data itself.

## Development

- [General development notes](development/README.md)
- [Agent Development](development/Agent.md).
- [Creating apps](development/Apps.md).
