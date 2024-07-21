#!/usr/bin/env bash

# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Stop on errors
set -e

# Install additional go build packages
go install github.com/magefile/mage@9e91a03eaa438d0d077aca5654c7757141536a60
