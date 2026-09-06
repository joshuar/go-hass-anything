# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

ARG ALPINE_VERSION=3.24.1@sha256:79ff19e9084a00eece421b2523fb93e22d730e2c0e525905de047e848e56d95f
ARG GO_VERSION=1.27.1-alpine3.24@sha256:f86f1a6701e3dcc445fec097a42f78b758f15950ccf032c2d3e54e2754d32fdb

FROM --platform=$BUILDPLATFORM docker.io/golang:${GO_VERSION} AS golang
FROM --platform=$BUILDPLATFORM docker.io/alpine:${ALPINE_VERSION} AS builder

COPY --from=golang /usr/local/go/ /usr/local/go/

ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"

# Import TARGETPLATFORM.
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Import APPDIR or set to examples directory.
ARG APPDIR=examples

# Install build requirements.
RUN apk add --update git linux-headers upx

# Set workdir.
WORKDIR /usr/src/go-hass-anything

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

# copy the src to the workdir
ADD . .

# remove unneeded dev/build directories
RUN rm -fr deployments dist/* || exit 0

# copy the user-specified APPDIR to a location that will be picked up during build
RUN rm -fr apps || exit 0
COPY $APPDIR apps/

# build the binary
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go tool mage -d build/magefiles -w . build:full

# compress binary with upx
RUN upx --best --lzma /usr/src/go-hass-anything/dist/go-hass-anything-$TARGETARCH*

FROM --platform=$BUILDPLATFORM alpine@sha256:79ff19e9084a00eece421b2523fb93e22d730e2c0e525905de047e848e56d95f

# Add image labels.
LABEL org.opencontainers.image.source="https://github.com/joshuar/go-hass-anything"
LABEL org.opencontainers.image.description="Send anything to Home Assistant, through MQTT, powered by Go"
LABEL org.opencontainers.image.licenses="MIT"

# Import TARGETPLATFORM and TARGETARCH
ARG TARGETPLATFORM
ARG TARGETARCH

# copy binary over from builder stage
COPY --from=builder /usr/src/go-hass-anything/dist/go-hass-anything-$TARGETARCH* /usr/bin/go-hass-anything

# allow custom uid and gid
ARG UID=1000
ARG GID=1000

# add user
RUN addgroup --gid "${GID}" go-hass-anything && \
    adduser --disabled-password --gecos "" --ingroup go-hass-anything \
    --uid "${UID}" go-hass-anything

# Set user.
USER go-hass-anything
ENTRYPOINT ["go-hass-anything"]
CMD ["run"]
