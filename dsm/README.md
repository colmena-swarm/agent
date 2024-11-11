# Distributed Service Manager (DSM)

The DSM is a small service that can start or stop roles on command, using Docker. 
It exposes an endpoint to receive service definitions passed to it by zenoh-client. It expects the Docker image ID to be part of the service definition.
Then it will wait for commands passed to it by the role-selector to endpoints `/start` and `/stop`

Build the service:
`go build`

Run the service:
`go run`

Run tests:
`go test ./...`