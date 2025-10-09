# syntax = docker/dockerfile:1-experimental

FROM --platform=$BUILDPLATFORM golang:1.24.2 AS build-stage
    ENV CGO_ENABLED=0
    COPY ./ /dsm
    WORKDIR /dsm
    RUN go mod download
    COPY . .
    ARG TARGETOS
    ARG TARGETARCH
    RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /dsmexec

FROM gcr.io/distroless/base-debian11 AS build-release-stage
    WORKDIR /
    COPY --from=build-stage /dsmexec /dsmexec

    # Run dsm as root
    # https://gitlab.bsc.es/wdc/projects/colmena/-/issues/12
    USER root

    # Run binary
    ENTRYPOINT ["/dsmexec"]