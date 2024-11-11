#!/usr/bin/env python3
#
#  Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#

DOCKER_COMPOSE_FILE="../compose-zenoh.yaml"

.PHONY: test integration-test clean zenoh-up zenoh-down

# Run all tests
test:
	cargo test --lib
	$(MAKE) integration-test

unit-test:
	cargo test --lib

# Run integration tests
integration-test: zenoh-up
	cargo test --test '*integration*'
	$(MAKE) zenoh-down

# Spin up Docker dependencies
zenoh-up:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Waiting for dependencies to start..."
	sleep 5  # Adjust or replace with health check if needed

# Shut down Docker dependencies
zenoh-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

# Clean up build artifacts
clean:
	cargo clean
