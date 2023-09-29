#!/usr/bin/env bash

# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Stop on errors
set -e

sudo apt -y update 

# Install mosquitto command-line utilities
export DEBIAN_FRONTEND=noninteractive && \
    sudo apt -y install mosquitto-clients micro

# Install additional go build packages
go install golang.org/x/tools/cmd/stringer@latest
go install golang.org/x/text/cmd/gotext@latest
go install github.com/matryer/moq@latest
go install github.com/goreleaser/goreleaser@latest
go get -u github.com/davecgh/go-spew/spew
