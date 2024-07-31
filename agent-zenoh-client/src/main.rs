/*
 *  Copyright 2002-2024 Barcelona Supercomputing Center (www.bsc.es)
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
use std::convert::TryFrom;
use std::env;
use std::str::FromStr;

use tonic::{transport::Server, Request, Response, Status};

use zenohclient::greeter_server::{Greeter, GreeterServer};
use zenohclient::{HelloReply, HelloRequest, MetricsQueryRequest, MetricsQueryResponse};
use zenoh::prelude::r#async::*;

pub mod zenohclient {
    tonic::include_proto!("zenohclient");
}

struct MyGreeter {
    session: Session
}

#[tonic::async_trait]
impl Greeter for MyGreeter {

    async fn say_hello(
        &self,
        request: Request<HelloRequest>,
    ) -> Result<Response<HelloReply>, Status> {
        println!("Got a request: {:?}", request);

        let reply = zenohclient::HelloReply {
            message: format!("Hello {}!", request.into_inner().name),
        };

        Ok(Response::new(reply))
    }

    async fn query_metrics(
        &self,
        request: Request<MetricsQueryRequest>,
    ) -> Result<Response<MetricsQueryResponse>, Status> {
        println!("Received request: {:?}", request);
        let inner = request.into_inner();
        let time_unit = inner.from_unit;
        let time_type = inner.from_type;
        let selector = inner.key + &format!("?_time=[now(-{}{})..]", time_unit, time_type);
        
        println!("Sending Query '{selector}'...");
        let mut values: Vec<f32> = Vec::new();
        let replies = self.session.get(selector).res().await.unwrap();
        while let Ok(reply) = replies.recv_async().await {
            match reply.sample {
                Ok(sample) => {
                    println!("Received ('{}': '{}')", sample.key_expr.as_str(), sample.value,);
                    let parsed = sample.value.to_string().parse::<f32>().unwrap();
                    values.push(parsed)
                },
                Err(err) => println!("Received (ERROR: '{}')", String::try_from(&err).unwrap()),
            }
        }
        if values.is_empty() {
            println!("No values retrieved from zenoh");
        }
        let kpi_met = compare_values_to_threshold(inner.comparison.as_str(), inner.threshold, values);
        println!("KPI met {}", kpi_met);
        let reply = zenohclient::MetricsQueryResponse {met: kpi_met};

        Ok(Response::new(reply))
    }
}

fn compare_values_to_threshold(comparison_str: &str, threshold: f32, values: Vec<f32>) -> bool{
    //https://gitlab.bsc.es/wdc/projects/colmena-group/agent/-/issues/16
    if values.is_empty() {
        return false;
    }

    let comparison = parse_comparison(comparison_str, threshold);
    println!("Comparing {} values to threshold: {}, comparison: {}", values.len(), threshold, comparison_str);
    return values.iter().all(comparison);
}

fn parse_comparison(comparison: &str, threshold: f32) -> Box<dyn Fn(&f32) -> bool> {
    match comparison {
        "<" => Box::new(move |x: &f32| *x < threshold),
        ">" => Box::new(move |x: &f32| *x > threshold),
        ">=" => Box::new(move |x: &f32| *x >= threshold),
        "<=" => Box::new(move |x: &f32| *x <= threshold),
        _   => panic!("Comparison not configured: {comparison}")
    }
}

async fn service_definition_subscription() {
    let peers = [EndPoint::from_str(&(env::var("ZENOH_ROUTER").unwrap())).unwrap()];
    let zenoh_client_config = config::client(peers);
    let zenoh_session = zenoh::open(zenoh_client_config).res().await.unwrap().into_arc();
    let key_expr = "colmena_service_definitions";
    let subscriber = zenoh_session.declare_subscriber(key_expr)
    .res()
    .await
    .unwrap();

    let client = reqwest::Client::new();
    println!("Starting Zenoh subscriber. key_expr: {:?}", key_expr);
    let _ = tokio::task::spawn(async move {
        while let Ok(sample) = subscriber.recv_async().await {
            let service_definition = sample.value.to_string();
            let parsed = json::parse(&service_definition).unwrap();
            let service_name = parsed["id"]["value"].as_str().unwrap();

            println!("Received service definition id: {:?}", service_name);

            let _ = client.post("http://agent:50551")
            .body(sample.value.to_string())
            .send()
            .await
            .and_then(|response| {println!("Forwarded service definition: {:?}, response: {:?}", service_name, response.status()); Ok(())})
            .or_else(|error| {println!("Error forwarding service definition: {:?}, error: {:?}", service_name, error); Err(error)});
        }
    }).await;
    println!("Finished Zenoh subscriber");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let server_address = "0.0.0.0:50051".parse().unwrap();
    
    //build gRPC service to query metrics
    let greeter = MyGreeter {
        session: zenoh::open(config::default()).res().await.unwrap()
    };

    tokio::spawn(async { service_definition_subscription().await });

    println!("starting Zenoh query gRPC service on {}", server_address);
    Server::builder()
        .add_service(GreeterServer::new(greeter))
        .serve(server_address)
        .await?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn given_a_comparison_with_no_values_higher_than_threshold_then_kpi_is_met() {
        let result = compare_values_to_threshold("<", 10.0, &vec![1.0, 2.0, 9.0]);
        assert_eq!(result, true);
    }

    #[test]
    fn given_a_comparison_with_a_value_equal_to_threshold_then_kpi_is_not_met() {
        let result = compare_values_to_threshold("<", 10.0, &vec![10.0, 2.0, 3.0]);
        assert_eq!(result, false);
    }

    #[test]
    fn given_a_comparison_with_no_values_then_kpi_is_not_met() {
        let lessthan = compare_values_to_threshold("<", 10.0, &vec![]);
        assert_eq!(lessthan, false);
        let greaterthan = compare_values_to_threshold(">", 10.0, &vec![]);
        assert_eq!(greaterthan, false);
    }

    #[test]
    fn equals_comparisons_work_as_expected() {
        let threshold = 10.0;
        let values = vec![10.0];
        let greater_than_or_equal_to_result = compare_values_to_threshold(">=", threshold, &values);
        let less_than_or_equal_to_result = compare_values_to_threshold("<=", threshold, &values);
        let less_than_result = compare_values_to_threshold("<", threshold, &values);
        let greater_than_result = compare_values_to_threshold(">", threshold, &values);
        assert_eq!(greater_than_or_equal_to_result, true);
        assert_eq!(less_than_or_equal_to_result, true);
        assert_eq!(less_than_result, false);
        assert_eq!(greater_than_result, false);
    }
}