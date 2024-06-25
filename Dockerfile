# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM golang@sha256:a8498215385dd85856145845f3caf34923fe5fbb11f3c7c1489ae43c4f263b20 AS builder

ARG APPDIR=pkg/apps

WORKDIR /usr/src/go-hass-anything

# copy the src to the workdir
ADD . .

# copy the user-specified APPDIR to a location that will be picked up during build
RUN rm -fr apps || exit 0
COPY $APPDIR apps/

# install mage
RUN go install github.com/magefile/mage@v1.15.0

# build the binary
RUN mage -v -d build/magefiles -w . build:full

FROM ubuntu@sha256:d11b1797008f48495a888a087b273f6581daef886da9d0bda9023557eff4f070

# import TARGETARCH
ARG TARGETARCH

# copy binary over from builder stage
COPY --from=builder /usr/src/go-hass-anything/dist/go-hass-anything-$TARGETARCH /usr/bin/go-hass-anything

ENTRYPOINT ["go-hass-anything"]
CMD ["run"]