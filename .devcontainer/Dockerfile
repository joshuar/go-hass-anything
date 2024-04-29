# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM mcr.microsoft.com/devcontainers/base:ubuntu
# Install the xz-utils package
ENV DEBIAN_FRONTEND noninteractive
RUN sudo apt -y update
RUN sudo apt -y install mosquitto-clients micro
RUN sudo rm -rf /var/lib/apt/lists/* /var/cache/apt/*
