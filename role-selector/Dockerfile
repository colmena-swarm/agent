# syntax = docker/dockerfile:1-experimental

FROM --platform=$BUILDPLATFORM golang:1.24.2 AS build-stage
    ENV CGO_ENABLED=0
    COPY ./ /role-selector
    WORKDIR /role-selector 
    RUN go mod download
    COPY . .
    ARG TARGETOS
    ARG TARGETARCH
    RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /roleselectorexec

FROM gcr.io/distroless/base-debian11 AS build-release-stage
    WORKDIR /
    COPY --from=build-stage /roleselectorexec /roleselectorexec
    ENTRYPOINT ["/roleselectorexec"]