# syntax=docker/dockerfile:1

# Build application from source
FROM golang:1.24.2 AS build-stage

    COPY ./ /role-selector
    WORKDIR /role-selector    
    RUN go mod download && \
        CGO_ENABLED=0 GOOS=linux go build -o /roleselectorexec


# Deploy application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
    WORKDIR /
    COPY --from=build-stage /roleselectorexec /roleselectorexec
    ENTRYPOINT ["/roleselectorexec"]