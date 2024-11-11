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

package roleselector

import (
	"context"
	"fmt"
	"log"
	"net"

	policies "colmena.bsc.es/role-selector/policy"
	policy "colmena.bsc.es/role-selector/policy/grpc"
	"colmena.bsc.es/role-selector/types"

	"google.golang.org/grpc"
)

type PolicyServer struct {
	policy.UnimplementedPolicyServiceServer
	initialized bool
	policyName  string
}

func NewPolicyServer() *PolicyServer {
	return &PolicyServer{}
}

func (s *PolicyServer) Initialize(ctx context.Context, req *policy.InitializeOrStopRequest) (*policy.InitializeOrStopResponse, error) {
	policyName := req.GetPolicyName()
	log.Printf("Initialize request received for policy: %s", policyName)
	s.initialized = true
	s.policyName = policyName
	return &policy.InitializeOrStopResponse{
		Success: true,
		Message: fmt.Sprintf("Policy '%s' initialized successfully", policyName),
	}, nil
}

func (s *PolicyServer) Decide(ctx context.Context, req *policy.DecideRequest) (*policy.DecideResponse, error) {
	log.Printf("Decide request received")

	var roles []*types.Role
	for _, role := range req.GetRoles() {
		roles = append(roles, &types.Role{Id: role.RoleName})
	}

	var kpis []types.KPI
	for _, lvl := range req.GetLevels() {
		kpis = append(kpis, types.KPI{
			Query:          lvl.GetName(),
			Value:          float64(lvl.GetValue()),
			Threshold:      float64(lvl.GetThreshold()),
			AssociatedRole: lvl.GetAssociatedRole(),
		})
	}

	var resources []types.Resource
	for _, res := range req.GetResources() {
		resources = append(resources, types.Resource{
			Name:  res.GetName(),
			Value: int(res.GetValue()),
		})
	}

	if s.policyName == "Eager" {
		decisions, err := (&policies.EagerPolicy{}).DecidePolicy(roles, kpis, resources)
		if err != nil {
			return nil, fmt.Errorf("policy decision failed: %w", err)
		}

		return &policy.DecideResponse{
			Decisions: decisions,
		}, nil
	}

	return nil, fmt.Errorf("unsupported policy name: %s", s.policyName)
}

func (s *PolicyServer) Stop(ctx context.Context, req *policy.InitializeOrStopRequest) (*policy.InitializeOrStopResponse, error) {
	log.Printf("Stop request received for policy: %s", req.GetPolicyName())
	s.initialized = false
	return &policy.InitializeOrStopResponse{
		Success: true,
		Message: fmt.Sprintf("Policy '%s' stopped successfully", req.GetPolicyName()),
	}, nil
}

func RunServer(address string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	policy.RegisterPolicyServiceServer(grpcServer, NewPolicyServer())

	go func() {
		log.Printf("gRPC server listening on %s", address)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return grpcServer, nil
}
