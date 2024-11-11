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
use std::time::Duration;


// zenoh constants
pub const ZENOH_KEY_EXPRESSION: &str = "colmena_service_definitions/*";
const DEFAULT_ZENOH_ROUTER: &str = "tcp/127.0.0.1:7447";
const ZENOH_SESSION_START_DELAY_KEY: &str = "ZENOH_SESSION_START_DELAY";
const ZENOH_SESSION_START_DELAY_DEFAULT: u64 = 10;

pub async fn get_published_service_definitions(downstream_services: Vec<String>) {
    // wait for zenoh session to start
    let zenoh_session_start_delay = env::var(ZENOH_SESSION_START_DELAY_KEY)
        .unwrap_or(ZENOH_SESSION_START_DELAY_DEFAULT.to_string());
    println!("Zenoh session start delay: {:?}", zenoh_session_start_delay);
    tokio::time::sleep(Duration::from_secs(zenoh_session_start_delay.parse::<u64>().unwrap())).await;

    println!("Getting published service definitions");
    let zenoh_session = zenoh_session().await;
    let replies = match zenoh_session.get(ZENOH_KEY_EXPRESSION).await {
        Ok(replies) => replies,
        Err(e) => {
            eprintln!("Error getting published service definitions from zenoh: {}", e);
            return;
        }
    };
    
    while let Ok(reply) = replies.recv_async().await {
        match reply.result() {
            Ok(result) => {
                let service_definition = result.payload().try_to_string().unwrap().into_owned();
                handle_service_definition(service_definition, downstream_services.clone()).await;
            }
            Err(e) => {
                eprintln!("Error getting service definition from zenoh replies: {}", e);
            }
        }
    }
}

pub async fn service_definition_subscription(downstream_services: Vec<String>) {
    let zenoh_session = zenoh_session().await;
    let subscriber = zenoh_session.declare_subscriber(ZENOH_KEY_EXPRESSION)
    .await.unwrap();

    let _ = tokio::task::spawn(async move {
        while let Ok(sample) = subscriber.recv_async().await {
            let service_definition = sample.payload().try_to_string().unwrap().into_owned();
            handle_service_definition(service_definition, downstream_services.clone()).await
        }
    }).await;
    println!("Zenoh subscriber stopped");
}

async fn handle_service_definition(service_definition: String, downstream_services: Vec<String>) {
    let parsed = json::parse(&service_definition).unwrap();
    let service_name = parsed["id"]["value"].as_str().unwrap().to_string();
    println!("Received service definition id: {:?}", service_name);

    for downstream_service in downstream_services.iter() {
        forward_service_definition(downstream_service, service_definition.clone()).await;
    }
}

async fn forward_service_definition(url: &str, service_definition: String) {
    let client = reqwest::Client::new();
    let parsed = json::parse(&service_definition).unwrap();
    let service_name = parsed["id"]["value"].as_str().unwrap().to_string();
    let _ = client.post(url)
    .header("Content-Type", "application/json")
    .body(service_definition)
    .send()
    .await
    .and_then(|response| {
        println!("Forwarded service definition: {:?}, recipient: {:?}, response: {:?}", service_name, url, response.status()); Ok(())})
    .or_else(|error| {
        println!("Error forwarding service definition: {:?}, recipient: {:?}, error: {:?}", service_name, url, error); Err(error)});
}

pub async fn zenoh_session() -> zenoh::Session {
    let mut zenoh_client_config = zenoh::Config::default();
    let zenoh_address = zenoh_address();
    
    let _ = zenoh_client_config.insert_json5(
        "connect/endpoints",
        &format!("[\"{}\"]", zenoh_address)
    );
    
    let _ = zenoh_client_config.insert_json5("mode", "client");
    zenoh::open(zenoh_client_config).await.unwrap()
}

pub fn zenoh_address() -> String {
    env::var("ZENOH_ROUTER").unwrap_or(DEFAULT_ZENOH_ROUTER.to_string())
}
