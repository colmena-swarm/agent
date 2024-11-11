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

package policies

import (
	"log"
	"testing"
	"time"

	"colmena.bsc.es/role-selector/policy/grpc/client"
	"colmena.bsc.es/role-selector/types"
)

func startTestServer(t *testing.T) func() {
	grpcServer, err := RunServer(":50051")
	if err != nil {
		t.Fatalf("Failed to start gRPC server: %v", err)
	}

	return func() {
		grpcServer.Stop()
		log.Println("Test gRPC server stopped")
	}
}

func TestInitializePolicy(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	time.Sleep(100 * time.Millisecond)

	client, err := client.NewGeneralClient("localhost:50051")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	resp, err := client.InitializePolicy("Eager")
	if err != nil {
		t.Fatalf("InitializePolicy failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("InitializePolicy returned failure: %s", resp.Message)
	}

	t.Log("InitializePolicy succeeded:", resp.Message)
}

func TestCallPolicy(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	time.Sleep(100 * time.Millisecond)

	client, err := client.NewGeneralClient("localhost:50051")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	roles := []types.DockerRoleDefinition{
		{Id: "example_role"},
	}

	levels := []types.KPI{
		{
			Query:          "cpu_utilization",
			Value:          10,
			Threshold:      5,
			AssociatedRole: "example_role",
		},
	}

	resources := []types.Resource{
		{Name: "cpu", Value: 10},
	}

	client.InitializePolicy("Eager")
	decisions, err := client.DecidePolicy(roles, levels, resources)
	if err != nil {
		t.Fatalf("DecidePolicy failed: %v", err)
	}

	t.Logf("Decisions: %v", decisions)
}

func TestStopPolicy(t *testing.T) {
	stop := startTestServer(t)
	defer stop()

	time.Sleep(100 * time.Millisecond)

	client, err := client.NewGeneralClient("localhost:50051")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	resp, err := client.StopPolicy("Eager")
	if err != nil {
		t.Fatalf("StopPolicy failed: %v", err)
	}
	if !resp.Success {
		t.Fatalf("StopPolicy returned failure: %s", resp.Message)
	}

	t.Log("StopPolicy succeeded:", resp.Message)
}
