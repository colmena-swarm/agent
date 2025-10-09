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
	"time"

	colmena_consensus "colmena.bsc.es/role-selector/policy/grpc/colmena_consensus"
	"colmena.bsc.es/role-selector/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GeneralClient struct {
	conn   *grpc.ClientConn
	client colmena_consensus.SelectionServiceClient
}

func NewGeneralClient(serverAddr string) (*GeneralClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := colmena_consensus.NewSelectionServiceClient(conn)

	return &GeneralClient{client: client}, nil
}

func (g *GeneralClient) Close() {
	if g.conn != nil {
		g.conn.Close()
	}
}

func (g *GeneralClient) RequestRoles(
	roles []*types.Role,
	kpis []types.KPI,
	resources []types.Resource,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, role := range roles {
		// Convert internal Role -> proto Role
		protoRole := &colmena_consensus.Role{
			RoleId:    role.Id,
			ServiceId: role.ServiceId,
			Resources: convertResources(role.Resources),
		}

		req := &colmena_consensus.RoleRequest{
			Role:        protoRole,
			StartOrStop: true,
		}

		_, err := g.client.RequestRoles(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to trigger role %s: %v", role.Id, err)
		}
	}

	return nil
}

// Helper to convert internal resources to proto resources
func convertResources(resources []types.Resource) []*colmena_consensus.Resource {
	protoResources := make([]*colmena_consensus.Resource, len(resources))
	for i, r := range resources {
		protoResources[i] = &colmena_consensus.Resource{
			Name:  r.Name,
			Value: r.Value,
		}
	}
	return protoResources
}
