<!--
 Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>

 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# General Development Notes

## Development Environment

It is recommended to use [Visual Studio Code](https://code.visualstudio.com/).
This project makes use of a [devcontainer](https://containers.dev/) to provide
some convenience during development.

If using Visual Studio Code, you should be prompted when opening your cloned
copy of the code to set up the dev container environment. The environment contains an
installation of Home Assistant and Mosquitto that can be used for testing. They should be
started automatically.

- Home Assistant will be listening on [localhost:8123](http://localhost:8123).
- Mosquitto will be listening on [localhost:1833](localhost:1833).

An example configuration for Mosquitto has been provided in
`deployments/mosquitto/config/mosquitto.conf.example`.

The Mosquitto command-line utilities (`mosquitto_{pub,sub}`) are installed in
the devcontainer.

## Building

Go Hass Anything makes use of `go generate` to generate some of the code. A
typical build process would be:

```shell
go generate ./...
go build
```

The resulting `go-hass-anything` binary in the root of the repo should be ready
to run.

## Packaging

Go Hass Anything uses [Goreleaser](https://goreleaser.com/intro/) to create
packages for Fedora, Arch, and Ubuntu.

To build a “local-only” package with Goreleaser:

```shell
goreleaser release --snapshot --clean
```

Packages will be available under the `dist/` folder. These can be installed
directly on your machine(s) using the typical package managers for your distribution.

See the [Goreleaser docs](https://goreleaser.com/quick-start/) for more commands
and information.

## Committing Code

This repository is using [conventional commit
messages](https://www.conventionalcommits.org/en/v1.0.0/#summary). This provides
the ability to automatically include relevant notes in the
[changelog](../../CHANGELOG.md). The [TL;DR](https://en.wikipedia.org/wiki/TL;DR)
is when writing commit messages, add a prefix:

- `feat:` for a new feature, like a new sensor.
- `fix:` when fixing an issue.
- `refactor:` when making non-visible but useful code changes.
- …and so on. See the link above or see the existing commit messages for examples.
