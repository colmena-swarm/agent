# syntax=docker/dockerfile:1

# Build application from source
FROM golang:1.22 AS build-stage

    COPY ./ /agent
    WORKDIR /agent    
    RUN go mod download && \
        CGO_ENABLED=0 GOOS=linux go build -o /agentexec


# Deploy application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
    WORKDIR /
    COPY --from=build-stage /agentexec /agentexec

    # Run agent as root
    # https://gitlab.bsc.es/wdc/projects/colmena/-/issues/12
    USER root

    # Run binary
    ENTRYPOINT ["/agentexec"]