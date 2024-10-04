# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM alpine@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0 AS builder
# Copy go from official image.
COPY --from=golang:1.23.2-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"
# Import TARGETPLATFORM.
ARG TARGETPLATFORM
# Import APPDIR or set to examples directory.
ARG APPDIR=examples
# Install build requirements.
RUN apk add --update git linux-headers
# Set workdir.
WORKDIR /usr/src/go-hass-anything
# copy the src to the workdir
ADD . .
# remove unneeded dev/build directories
RUN rm -fr deployments dist/* || exit 0
# copy the user-specified APPDIR to a location that will be picked up during build
RUN rm -fr apps || exit 0
COPY $APPDIR apps/
# install mage
RUN go install github.com/magefile/mage@9e91a03eaa438d0d077aca5654c7757141536a60 # v1.15.0
# build the binary
RUN mage -v -d build/magefiles -w . build:full

FROM alpine@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0
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
