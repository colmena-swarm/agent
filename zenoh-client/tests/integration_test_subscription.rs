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

#[cfg(test)]
mod tests {
    use zenoh_client::{service_definition_subscription, zenoh_session};
    
    #[tokio::test(flavor = "multi_thread", worker_threads = 5)]
    async fn test_subscription_forwards_service_definition_correctly() {
        let service_definition = json::parse(r#"{ "id": { "value": "test_service" } }"#).unwrap();

        let server = httpmock::MockServer::start();
        let mock_endpoint = server.mock(|when, then| {
            when.method(httpmock::Method::POST)
                .path("/xyz");
            then.status(200);
        });

        let downstream_services = vec![server.url("/xyz")];
        let handle = tokio::spawn(async { service_definition_subscription(downstream_services).await });

        tokio::time::sleep(std::time::Duration::from_secs(1)).await;

        let session = zenoh_session().await;
        session.put("colmena_service_definitions/expression", service_definition.to_string().as_bytes()).await.unwrap();

        tokio::time::sleep(std::time::Duration::from_secs(1)).await;

        mock_endpoint.assert();

        session.close().await.unwrap();
        handle.abort();
    }
}
