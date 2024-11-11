# Zenoh Client
Zenoh clients don't exist for the Go programming language as of May 2025. This small Rust app exists to subscribe to COLMENA services that are published onto the Zenoh network and then use HTTP to publish the service to components within COLMENA.

Build the service locally:
`cargo build`

Run the service locally:
`cargo run`

Docker image is built using https://github.com/LukeMathWalker/cargo-chef

Environment variables:
| Environment variable      | Explanation                                                       |
| -------------             | -------------                                                     |
| ZENOH_ROUTER              | hostport of the Zenoh router defaults to tcp/172.0.0.1:7447       |
| DOWNSTREAM_SERVICES       | downstream services which will receive the service description    |
| ZENOH_SESSION_START_DELAY | delay before getting services from zenoh, while docker compose is starting   |



