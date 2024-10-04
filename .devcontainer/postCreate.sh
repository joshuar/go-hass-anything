#!/usr/bin/env bash

# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Stop on errors
set -e

# Install additional go build packages
go install github.com/magefile/mage@9e91a03eaa438d0d077aca5654c7757141536a60
go install github.com/sigstore/cosign/v2/cmd/cosign@fb651b4ddd8176bd81756fca2d988dd8611f514d

# Install and configure starship
curl -sS https://starship.rs/install.sh | sh -s -- -y || exit -1
mkdir -p ~/.config/fish
echo "starship init fish | source" >> ~/.config/fish/config.fish
exit 0
