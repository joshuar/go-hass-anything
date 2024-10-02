#!/usr/bin/env bash

# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Stop on errors
set -e

# Install additional go build packages
go install github.com/magefile/mage@9e91a03eaa438d0d077aca5654c7757141536a60
go install golang.org/x/tools/cmd/stringer@v0.23.0
go install golang.org/x/text/cmd/gotext@v0.16.0
go install github.com/matryer/moq@0bf2e8a069abaefdfd07e4902d204441cca17298

# Install go tools
go install github.com/tomwright/dasel/v2/cmd/dasel@5d94a3049c2e81a410a6f48cf084c86c98393797
go install github.com/sigstore/cosign/v2/cmd/cosign@fb651b4ddd8176bd81756fca2d988dd8611f514d

# Install and configure starship
curl -sS https://starship.rs/install.sh | sh -s -- -y || exit -1
mkdir -p ~/.config/fish
echo "starship init fish | source" >> ~/.config/fish/config.fish
exit 0
