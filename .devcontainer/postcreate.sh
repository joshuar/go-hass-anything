#!/usr/bin/bash
# Copyright 2025 Joshua Rich <joshua.rich@gmail.com>.
# SPDX-License-Identifier: MIT

set -x

cd /workspace

# Install Go tools.
go mod tidy
go install github.com/air-verse/air@latest
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sigstore/cosign/v3/cmd/cosign@latest

exit 0
