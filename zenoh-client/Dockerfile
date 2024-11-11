FROM lukemathwalker/cargo-chef:latest-rust-1.83.0 AS chef
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
  RUN cargo build --release

# Without Rust toolchain
FROM debian:bookworm-slim AS runtime
  WORKDIR /app
  COPY --from=builder /app/target/release/zenoh-client /usr/local/bin
  RUN apt-get update && apt install -y openssl
  ENTRYPOINT ["/usr/local/bin/zenoh-client"]
  