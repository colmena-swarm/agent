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
	"log"
	"os"
	"sync"
	"time"

	"colmena.bsc.es/role-selector/policy"
	"colmena.bsc.es/role-selector/servicedescription"
	"colmena.bsc.es/role-selector/sla"
	"colmena.bsc.es/role-selector/types"
)

type RoleSelector struct {
	ServiceDescriptionChan chan types.ServiceDescription
	AlertsChan             chan types.Alert
	Hardware               string
	Policy                 policy.Policy
	RoleRunner             RoleRunner
	KpiRetriever           sla.KpiRetriever
	RolesByServiceId       map[string][]*types.Role

	mu sync.Mutex
}

// Default for now, in the future take/calculate somewhere from role execution
var DefaultResources = []types.Resource{
	{Name: "core", Value: 30},
	{Name: "ram", Value: 30},
	{Name: "disk", Value: 30},
}

const defaultTickerInterval = 10 * time.Second

func tickerInterval() time.Duration {
	interval := os.Getenv("ROLE_SELECTION_INTERVAL")
	if interval == "" {
		log.Printf("Role selection interval not set, using default ticker interval: %v", defaultTickerInterval)
		return defaultTickerInterval
	}
	duration, err := time.ParseDuration(interval)
	if err != nil {
		log.Printf("Invalid role selection interval: %v, using default ticker interval: %v", err, defaultTickerInterval)
		return defaultTickerInterval
	}
	log.Printf("Role selection ticker interval: %v", duration)
	return duration
}

func (rs *RoleSelector) SetPolicy(policy policy.Policy) {
	rs.Policy = policy
}

func (rs *RoleSelector) GetRoleByName(serviceId string, roleName string) (*types.Role, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	roles, ok := rs.RolesByServiceId[serviceId]
	if !ok {
		return nil, false
	}

	for _, role := range roles {
		if role.Id == roleName {
			return role, true
		}
	}
	return nil, false
}

func (rs *RoleSelector) Run(stoppedRoles chan StoppedEvent) {
	rs.RolesByServiceId = make(map[string][]*types.Role)
	ticker := time.NewTicker(tickerInterval())
	defer ticker.Stop()
	for {
		select {
		case serviceDescription := <-rs.ServiceDescriptionChan:
			serviceId := serviceDescription.ServiceId.Value
			log.Println("Received service description for service: ", serviceId)
			_, roleDescriptions, _ := servicedescription.Parse(&serviceDescription, rs.Hardware)
			if len(roleDescriptions) == 0 {
				log.Printf("No roles matching hardware requirements %v found for service: %v", rs.Hardware, serviceId)
				continue
			}

			roles := parse(roleDescriptions, serviceId)

			/* handle updating a service
			this puts the roles into an updating state to prevent the role selector from starting
			either version of the role until the previous has been stopped */
			if _, ok := rs.RolesByServiceId[serviceId]; ok {
				log.Printf("Updating roles for service: %v", serviceId)
				Clean(serviceId, rs.RolesByServiceId[serviceId], roles, rs.RoleRunner)
			}

			rs.RolesByServiceId[serviceId] = roles
			rs.checkRoles()

		case alert := <-rs.AlertsChan:
			serviceId := alert.ServiceId
			roles := rs.RolesByServiceId[serviceId]
			kpis := types.ParseKpis(alert.Kpis)

			rs.decide(roles, kpis, serviceId)

		//case <-ticker.C:
		//periodically check all services
		//	rs.checkRoles()

		case stopped := <-stoppedRoles:
			rs.handleRoleStopped(stopped)
		}
	}
}

func (rs *RoleSelector) checkRoles() {
	for serviceId, roles := range rs.RolesByServiceId {
		kpis, err := rs.KpiRetriever.Get(serviceId)
		if err != nil {
			log.Printf("Failed to get KPIs for service %v: %v", serviceId, err)
			continue
		}
		rs.decide(roles, kpis, serviceId)
	}
}

func (rs *RoleSelector) runMandatoryRoles(roles []*types.Role) []*types.Role {
	decisions := make(map[string]bool)
	var remainingRoles []*types.Role
	var mandatoryRoles []*types.Role

	for _, role := range roles {
		if len(role.Kpis) == 0 {
			decisions[role.Id] = true // mark mandatory roles as running
			mandatoryRoles = append(mandatoryRoles, role)
		} else {
			remainingRoles = append(remainingRoles, role) // keep roles with KPIs
		}
	}

	rs.executeDecisions(decisions, mandatoryRoles)
	return remainingRoles
}

func (rs *RoleSelector) decide(roles []*types.Role, kpis []types.KPI, serviceId string) {
	remainingRoles := rs.runMandatoryRoles(roles)
	decisions, _ := rs.Policy.DecidePolicy(remainingRoles, kpis, []types.Resource{})
	if len(decisions) > 0 {
		rs.executeDecisions(decisions, remainingRoles)
	}
}

func (rs *RoleSelector) executeDecisions(decisions map[string]bool, roles []*types.Role) {
	for _, currentRole := range roles {
		decision := decisions[currentRole.Id]
		rs.executeDecision(currentRole, decision)
	}
}

func (rs *RoleSelector) TriggerDecision(decision types.Decision) {
	log.Printf("Received decision %v %v %v", decision.RoleId, decision.ServiceId, decision.StartOrStop)
	role, _ := rs.GetRoleByName(decision.ServiceId, decision.RoleId)
	rs.executeDecision(role, decision.StartOrStop)
}

func (rs *RoleSelector) executeDecision(role *types.Role, decision bool) {
	switch {
	case decision && role.State == types.Stopped:
		log.Printf("Decided to run role: %v ...", role.Id)
		rs.RoleRunner.Run(role.Id, role.ServiceId, role.ImageId)
		role.State = types.Running

	case !decision && role.State == types.Running:
		log.Printf("Decided to stop role: %v ...", role.Id)
		rs.RoleRunner.Stop(role.Id, role.ServiceId, role.ImageId, false)
		role.State = types.Stopped
	}
}

func (rs *RoleSelector) handleRoleStopped(stopped StoppedEvent) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	log.Printf("Received stopped event for role: %v, serviceId: %v", stopped.RoleId, stopped.ServiceId)
	roles := rs.RolesByServiceId[stopped.ServiceId]
	for _, role := range roles {
		if role.Id == stopped.RoleId {
			switch role.State {
			case types.Updating:
				log.Printf("Previous version stopped role: %v, serviceId: %v, imageId: %v", stopped.RoleId, stopped.ServiceId, stopped.ImageId)
				role.State = types.Stopped
				rs.checkRoles()
			case types.Running:
				log.Printf("Role stopped roleId: %v, serviceId: %v, imageId: %v", stopped.RoleId, stopped.ServiceId, stopped.ImageId)
				role.State = types.Stopped
			case types.Stopped:
				log.Printf("Error: Role already stopped roleId: %v, serviceId: %v, imageId: %v", stopped.RoleId, stopped.ServiceId, stopped.ImageId)
				role.State = types.Stopped
			}
		}
	}
}

func (rs *RoleSelector) GetRoleStatus(roleId string, serviceId string) (types.RoleState, bool) {
	roles := rs.RolesByServiceId[serviceId]
	for _, role := range roles {
		if role.Id == roleId {
			return role.State, true
		}
	}
	return types.Unknown, false
}
