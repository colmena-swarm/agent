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

package policy

import (
	"context"
	"fmt"
	"log"
	"time"

	"colmena.bsc.es/role-selector/policy/grpc/client"
	"colmena.bsc.es/role-selector/policy/grpc/server"
	"colmena.bsc.es/role-selector/types"
)

// ExecutionMode represents whether a policy is synchronous or asynchronous
type ExecutionMode int

const (
	Synchronous ExecutionMode = iota
	Asynchronous
)

type Policy interface {
	DecidePolicy(
		roles []*types.Role,
		kpis []types.KPI,
		resources []types.Resource,
	) (map[string]bool, error)
	Stop()
	Name() string
	Mode() ExecutionMode
}

type LazyPolicy struct{}

func (p *LazyPolicy) Name() string {
	return "Lazy"
}

func (p *LazyPolicy) Mode() ExecutionMode {
	return Synchronous
}

func (p *LazyPolicy) DecidePolicy(
	roles []*types.Role,
	kpis []types.KPI,
	resources []types.Resource,
) (map[string]bool, error) {
	decisions := make(map[string]bool)
	roleSet := make(map[string]bool)

	for _, role := range roles {
		roleSet[role.Id] = (role.State == types.Running || role.State == types.Updating)
	}

	// Process KPIs
	for _, kpi := range kpis {
		log.Printf("KPI query: %v, associated role: %v, level: %v", kpi.Query, kpi.AssociatedRole, kpi.Level)
		if kpi.AssociatedRole != "" && (kpi.Level == "Broken" || kpi.Level == "Critical") {
			if _, exists := roleSet[kpi.AssociatedRole]; exists {
				roleSet[kpi.AssociatedRole] = true
			}
		} else {
			if _, exists := roleSet[kpi.AssociatedRole]; exists {
				roleSet[kpi.AssociatedRole] = false
			}
		}
	}

	for roleID, decision := range roleSet {
		decisions[roleID] = decision
	}

	return decisions, nil
}

func (p *LazyPolicy) Stop() {}

type EagerPolicy struct{}

func (p *EagerPolicy) Name() string {
	return "Eager"
}

func (p *EagerPolicy) Mode() ExecutionMode {
	return Synchronous
}

func (p *EagerPolicy) DecidePolicy(
	roles []*types.Role,
	kpis []types.KPI,
	resources []types.Resource,
) (map[string]bool, error) {
	decisions := make(map[string]bool)
	for _, role := range roles {
		decisions[role.Id] = true
	}
	return decisions, nil
}

func (p *EagerPolicy) Stop() {}

type KPIMetEvent struct {
	RoleID  string
	KPI     string
	Elapsed float64 // seconds
}

type ConsensusPolicy struct {
	Client          *client.GeneralClient
	Server          *server.Server
	DecisionTrigger func(types.Decision)
	cancel          context.CancelFunc

	// Per-role tracking
	LastCallPerRole map[string]time.Time     // roleID -> last decision request time
	CooldownPerRole map[string]time.Duration // roleID -> cooldown duration
}

func (p *ConsensusPolicy) Name() string {
	return "Consensus"
}

func (p *ConsensusPolicy) Mode() ExecutionMode {
	return Asynchronous
}

// NewConsensusPolicy creates a new ConsensusPolicy with a gRPC client and server
func NewConsensusPolicy(endpoint string, trigger func(types.Decision)) (*ConsensusPolicy, error) {
	client, err := client.NewGeneralClient(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create consensus client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := server.NewServer(":50055", trigger)
	go server.Start(ctx)

	return &ConsensusPolicy{
		Client:          client,
		Server:          server,
		DecisionTrigger: trigger,
		cancel:          cancel,
		LastCallPerRole: make(map[string]time.Time),
		CooldownPerRole: make(map[string]time.Duration),
	}, nil
}

// DecidePolicy decides which roles to request
func (p *ConsensusPolicy) DecidePolicy(
	roles []*types.Role,
	kpis []types.KPI,
	resources []types.Resource,
) (map[string]bool, error) {

	now := time.Now()
	roleSet := make(map[string]struct{})
	decisions := make(map[string]bool)

	// Track roles with broken/critical KPIs
	for _, kpi := range kpis {
		if kpi.AssociatedRole == "" {
			continue
		}

		// KPI Broken/Critical → role needs request
		if kpi.Level == "Broken" || kpi.Level == "Critical" {
			roleSet[kpi.AssociatedRole] = struct{}{}

			// Check cooldown per role
			lastCall, exists := p.LastCallPerRole[kpi.AssociatedRole]
			cooldown := p.CooldownPerRole[kpi.AssociatedRole]
			if !exists {
				cooldown = 10 * time.Minute // default
				p.CooldownPerRole[kpi.AssociatedRole] = cooldown
			}
			if exists && time.Since(lastCall) < cooldown {
				remaining := cooldown - time.Since(lastCall)
				log.Printf("[ConsensusPolicy] Role %s cooldown active: %.1f seconds remaining, skipping request",
					kpi.AssociatedRole, remaining.Seconds())
				delete(roleSet, kpi.AssociatedRole) // skip this role for now
			}

		} else {
			// KPI Met -> log time since last request
			if lastCall, ok := p.LastCallPerRole[kpi.AssociatedRole]; ok {
				elapsed := time.Since(lastCall)
				log.Printf("[ConsensusPolicy] KPI Met for role %s, elapsed since last request: %.1f seconds",
					kpi.AssociatedRole, elapsed.Seconds())
			}
			//delete(p.LastCallPerRole, kpi.AssociatedRole)
			decisions[kpi.AssociatedRole] = false
		}

	}

	// Filter roles to request
	var rolesToRequest []*types.Role
	for _, role := range roles {
		if _, ok := roleSet[role.Id]; ok {
			rolesToRequest = append(rolesToRequest, role)
		}
	}

	if len(rolesToRequest) == 0 {
		log.Printf("[ConsensusPolicy] No roles need requesting, skipping gRPC call")
		return decisions, nil
	}

	// Send request
	log.Printf("[ConsensusPolicy] Requesting roles for %d affected role(s)", len(rolesToRequest))
	if err := p.Client.RequestRoles(rolesToRequest, kpis, resources); err != nil {
		log.Printf("[ConsensusPolicy] Error requesting roles: %v", err)
		return decisions, err
	}

	// Update last call timestamp per role
	for _, role := range rolesToRequest {
		p.LastCallPerRole[role.Id] = now
	}

	log.Printf("[ConsensusPolicy] Decision request sent at %s", now.Format(time.RFC3339))

	return decisions, nil
}

func (p *ConsensusPolicy) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	if p.Server != nil {
		p.Server.Stop()
	}
}
