# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM golang:1.21

WORKDIR /usr/src/go-hass-anything

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go generate ./...
RUN go build -v -o /usr/local/bin/go-hass-anything ./...

CMD ["go-hass-anything", "run"]