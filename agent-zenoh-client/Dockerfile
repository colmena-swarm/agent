FROM lukemathwalker/cargo-chef:0.1.64-rust-latest AS chef
WORKDIR /app

FROM chef AS planner
  COPY . .
  RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder 
  COPY --from=planner /app/recipe.json recipe.json
  # Build dependencies - this is the caching Docker layer!
  RUN cargo chef cook --release --recipe-path recipe.json
  # Build application
  COPY . .
  RUN apt-get update && \
      DEBIAN_FRONTEND=noninteractive \
      apt-get install --no-install-recommends --assume-yes \
        protobuf-compiler \
      && \
      cargo build --release

# Without Rust toolchain
FROM debian:bookworm-slim AS runtime
  WORKDIR /app
  COPY --from=builder /app/target/release/agent-zenoh /usr/local/bin
  RUN apt-get update && apt install -y openssl
  ENTRYPOINT ["/usr/local/bin/agent-zenoh"]
  