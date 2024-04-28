# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM docker.io/library/golang:1.22
ARG APPDIR

WORKDIR /usr/src/go-hass-anything

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# copy the user-specified APPDIR to a location that will be picked up during build
RUN test -L apps && rm apps || exit 0
COPY $APPDIR apps/

RUN go install github.com/matryer/moq@latest
RUN go install golang.org/x/tools/cmd/stringer@latest
RUN go generate ./...
RUN go build -v -o /go/bin/go-hass-anything

RUN useradd -ms /bin/bash gouser
USER gouser
WORKDIR /home/gouser

ENTRYPOINT ["go-hass-anything"]
CMD ["run"]