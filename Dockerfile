# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM --platform=$BUILDPLATFORM alpine AS builder

ARG TARGETARCH
ARG APPDIR=pkg/apps

RUN apk add --update go git

ENV PATH="$PATH:/root/go/bin"

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

FROM --platform=$BUILDPLATFORM alpine

# import TARGETARCH
ARG TARGETARCH

# copy binary over from builder stage
COPY --from=builder /usr/src/go-hass-anything/dist/go-hass-anything-$TARGETARCH /usr/bin/go-hass-anything

ENTRYPOINT ["go-hass-anything"]
CMD ["run"]