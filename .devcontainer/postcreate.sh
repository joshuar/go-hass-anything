#!/usr/bin/bash
# Copyright 2025 Joshua Rich <joshua.rich@gmail.com>.
# SPDX-License-Identifier: MIT

set -x

cd /workspace


# Install Go packages.
echo 'set --export PATH "$HOME/go/bin" /go/bin /usr/local/go/bin $PATH' >> ~/.config/fish/config.fish
echo 'export PATH="$HOME/go/bin:/go/bin:/usr/local/go/bin:$PATH' >> ~/.bashrc
if [[ -e go.mod ]]; then
    export PATH="$HOME/go/bin:/go/bin:/usr/local/go/bin:$PATH" && \
        go mod tidy && \
        go install golang.org/x/tools/gopls@v0.23.0 && \
        go install github.com/air-verse/air@v1.67.4 && \
        go install github.com/magefile/mage@v1.17.2 && \
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/refs/tags/v2.13.2/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.13.2
fi

if [[ -e ./.custom-gcl.yaml ]]; then
    golangci-lint custom && \
    mv /tmp/golangci-lint-v2 $(go env GOPATH)/bin/
fi

# Install oh-my-posh.
mkdir -p ~/.local/bin && curl -s https://raw.githubusercontent.com/JanDeDobbeleer/oh-my-posh/refs/tags/v31.1.2/website/static/install.sh | bash -s
# Set up shells to use oh-my-posh.
mkdir -p ~/.config/fish \
    && echo "~/.local/bin/oh-my-posh init fish | source" >>~/.config/fish/config.fish \
    && echo 'eval "$(~/.local/bin/oh-my-posh init bash)""' >>~/.bashrc

exit 0
