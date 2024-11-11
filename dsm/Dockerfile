# syntax=docker/dockerfile:1

# Build application from source
FROM golang:1.24.2 AS build-stage

    COPY ./ /dsm
    WORKDIR /dsm    
    RUN go mod download && \
        CGO_ENABLED=0 GOOS=linux go build -o /dsmexec


# Deploy application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
    WORKDIR /
    COPY --from=build-stage /dsmexec /dsmexec

    # Run dsm as root
    # https://gitlab.bsc.es/wdc/projects/colmena/-/issues/12
    USER root

    # Run binary
    ENTRYPOINT ["/dsmexec"]