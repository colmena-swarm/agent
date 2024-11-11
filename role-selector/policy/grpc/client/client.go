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

package client

import (
	"context"
	"fmt"

	policy "colmena.bsc.es/role-selector/policy/grpc"
	"colmena.bsc.es/role-selector/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GeneralClient struct {
	conn   *grpc.ClientConn
	client policy.PolicyServiceClient
}

func NewGeneralClient(serverAddr string) (*GeneralClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := policy.NewPolicyServiceClient(conn)

	return &GeneralClient{client: client}, nil
}

func (g *GeneralClient) Close() {
	if g.conn != nil {
		g.conn.Close()
	}
}

func (g *GeneralClient) InitializePolicy(policyName string) (*policy.InitializeOrStopResponse, error) {
	return g.client.Initialize(context.Background(), &policy.InitializeOrStopRequest{PolicyName: policyName})
}

func (g *GeneralClient) DecidePolicy(
	roles []types.DockerRoleDefinition,
	levels []types.KPI,
	resources []types.Resource,
) (map[string]bool, error) {
	var levelPtrs []*policy.IndicatorLevel
	var resourcePtrs []*policy.Resource
	var rolePtrs []*policy.Role

	for _, l := range levels {
		levelPtrs = append(levelPtrs, &policy.IndicatorLevel{
			Name:           l.Query,
			Value:          float32(l.Value),
			Threshold:      float32(l.Threshold),
			AssociatedRole: l.AssociatedRole,
		})
	}

	for _, r := range resources {
		resourcePtrs = append(resourcePtrs, &policy.Resource{
			Name:  r.Name,
			Value: float32(r.Value),
		})
	}

	// Hardcoded for now
	roleResources := []*policy.Resource{
		{
			Name:  "core",
			Value: 0.3,
		},
		{
			Name:  "disk",
			Value: 0.3,
		},
		{
			Name:  "ram",
			Value: 0.3,
		},
	}

	for _, role := range roles {
		rolePtrs = append(rolePtrs, &policy.Role{
			RoleName:  role.Id,
			IsRunning: false,
			Resources: roleResources,
		})
	}

	decideResponse, err := g.client.Decide(context.Background(), &policy.DecideRequest{
		Roles:     rolePtrs,
		Levels:    levelPtrs,
		Resources: resourcePtrs,
	})

	return decideResponse.Decisions, err
}

func (g *GeneralClient) StopPolicy(policyName string) (*policy.InitializeOrStopResponse, error) {
	return g.client.Stop(context.Background(), &policy.InitializeOrStopRequest{PolicyName: policyName})
}
