/*
 *  Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

use std::env;
use zenoh_client::{get_published_service_definitions, service_definition_subscription, zenoh_address};

fn downstream_services() -> Vec<String> {
    env::var("DOWNSTREAM_SERVICES")
        .unwrap_or_default()
        .split(',')
        .map(|s| s.trim().to_string())
        .filter(|s| !s.is_empty())
        .collect()
}

#[tokio::main]
async fn main() {
    println!("Zenoh address: {:?}", zenoh_address());
    get_published_service_definitions(downstream_services()).await;
    println!("Finished getting published service definitions, starting subscription...");
    let handle = tokio::spawn(async { service_definition_subscription(downstream_services()).await });
    let _ = handle.await;
}
