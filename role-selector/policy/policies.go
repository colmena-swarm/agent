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
	"log"

	"colmena.bsc.es/role-selector/types"
)

type Policy interface {
	DecidePolicy(
		roles []*types.Role,
		levels []types.KPI,
		resources []types.Resource,
	) (map[string]bool, error)
	Name() string
}

type LazyPolicy struct{}

func (p *LazyPolicy) Name() string {
	return "Lazy"
}

type EagerPolicy struct{}

func (p *EagerPolicy) Name() string {
	return "Eager"
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

func (p *LazyPolicy) DecidePolicy(
	roles []*types.Role,
	kpis []types.KPI,
	resources []types.Resource,
) (map[string]bool, error) {

	decisions := make(map[string]bool)

	roleSet := make(map[string]bool)
	for _, role := range roles {
		roleSet[role.Id] = role.IsRunning
	}

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
