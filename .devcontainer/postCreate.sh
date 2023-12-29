#!/usr/bin/env bash

# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Stop on errors
set -e

cd /workspaces/go-hass-anything && go mod tidy

# Install additional go build packages
go install golang.org/x/tools/cmd/stringer@latest
go install golang.org/x/text/cmd/gotext@latest
go install github.com/matryer/moq@latest
go install github.com/goreleaser/goreleaser@latest
go get -u github.com/davecgh/go-spew/spew
