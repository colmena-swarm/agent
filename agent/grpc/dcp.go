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
package grpc

import (
	context "context"
	"flag"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	addr = flag.String("addr", "127.0.0.1:5555", "dcp address")
)

func GetAllServices(found chan []*DockerRoleDefinition) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewColmenaPlatformClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.GetAllServices(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get all services: %v", err)
	}
	parseServiceDescriptionStream(stream, found)
}

func SubscribeToServices(found chan []*DockerRoleDefinition) {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewColmenaPlatformClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := c.SubscribeToServiceChanges(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not subscribe to service changes: %v", err)
		return
	}
	parseServiceDescriptionStream(stream, found)
}

func parseServiceDescriptionStream(stream ColmenaPlatform_GetAllServicesClient, found chan []*DockerRoleDefinition) {
	for {
		serviceDescription, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to get services: %v", err)
		}
		found <- serviceDescription.DockerRoleDefinitions
	}
}
