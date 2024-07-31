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
package zenohclient

import (
	context "context"
	"flag"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	zenohaddr = flag.String("zenohaddr", "zenoh-client:50051", "the address to connect to")
)


// uses zenoh connector service. more info: https://gitlab.bsc.es/wdc/projects/colmena/-/issues/11
func MetricsMet(request *MetricsQueryRequest) bool {
	flag.Parse()
	conn, err := grpc.NewClient(*zenohaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to zenoh: %v", err)
	}
	defer conn.Close()

	c := NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Print(request)

	r, err := c.QueryMetrics(ctx, request)
	if err != nil {
		log.Print(err.Error())
		log.Fatalf("could not query metrics: %v", err)
	}
	return r.Met
}