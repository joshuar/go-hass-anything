# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM mcr.microsoft.com/devcontainers/base:ubuntu
# Install the xz-utils package
ARG DEBIAN_FRONTEND noninteractive
RUN sudo apt-get -y update
RUN sudo apt-get -y install mosquitto-clients micro fish
RUN sudo rm -rf /var/lib/apt/lists/* /var/cache/apt/*
