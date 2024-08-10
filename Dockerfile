# Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

FROM --platform=$BUILDPLATFORM alpine@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0 AS builder

ARG TARGETARCH
ARG APPDIR=examples

RUN apk add --update go git linux-headers

ENV PATH="$PATH:/root/go/bin"

WORKDIR /usr/src/go-hass-anything

# copy the src to the workdir
ADD . .

# remove unneeded dev/build directories
RUN rm -fr deployments dist/* || exit 0

# copy the user-specified APPDIR to a location that will be picked up during build
RUN rm -fr apps || exit 0
COPY $APPDIR apps/

# install mage
RUN go install github.com/magefile/mage@v1.15.0

# build the binary
RUN mage -v -d build/magefiles -w . build:full

FROM --platform=$BUILDPLATFORM alpine@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0

# allow custom uid and gid
ARG UID=1000
ARG GID=1000

# add user
RUN addgroup --gid "${GID}" go-hass-anything && \
    adduser --disabled-password --gecos "" --ingroup go-hass-anything \
        --uid "${UID}" go-hass-anything

# import TARGETARCH
ARG TARGETARCH

# copy binary over from builder stage
COPY --from=builder /usr/src/go-hass-anything/dist/go-hass-anything-$TARGETARCH /usr/bin/go-hass-anything

USER go-hass-anything

ENTRYPOINT ["go-hass-anything"]
CMD ["run"]