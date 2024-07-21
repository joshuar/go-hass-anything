# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM mcr.microsoft.com/devcontainers/base:ubuntu@sha256:08a925b41755bcff957285d31e9bf5f0b7fc74962f455a929c15dc341b984b6f
# Install the xz-utils package
ARG DEBIAN_FRONTEND noninteractive
RUN sudo apt-get -y update
RUN sudo apt-get -y install mosquitto-clients micro fish ripgrep fzf
RUN sudo rm -rf /var/lib/apt/lists/* /var/cache/apt/*
