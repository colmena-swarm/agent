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

func (rs *RoleSelector) Run() {
	rolesByServiceId := make(map[string][]*types.Role)
	ticker := time.NewTicker(tickerInterval())
	defer ticker.Stop()
	for {
		select {
		case serviceDescription := <-rs.ServiceDescriptionChan:
			serviceId := serviceDescription.ServiceId.Value
			log.Println("Received service description for service: ", serviceId)
			_, roleDescriptions, _ := servicedescription.Parse(&serviceDescription, rs.Hardware)
			if len(roleDescriptions) == 0 {
				log.Printf("No roles found for service: %v", serviceId)
				continue
			}
			var roles []*types.Role
			for _, roleDescription := range roleDescriptions {
				toRun := false
				if len(roleDescription.Kpis) == 0 {
					// Run roles that don't have an associated KPI.
					log.Printf("isRunning for role %v to true", roleDescription.Id)
					toRun = true
					log.Printf("Decided to run role: %v, serviceId: %v, imageId: %v", roleDescription.Id, serviceId, roleDescription.ImageId)
					rs.RoleRunner.Run(roleDescription.Id, serviceId, roleDescription.ImageId)
				}
				role := &types.Role{
					Id:        roleDescription.Id,
					ImageId:   roleDescription.ImageId,
					IsRunning: toRun,
					Resources: DefaultResources,
				}
				roles = append(roles, role)
			}

			rolesByServiceId[serviceId] = roles

		case alert := <-rs.AlertsChan:
			serviceId := alert.ServiceId
			log.Println("Received alert for service: ", serviceId)
			roles := rolesByServiceId[serviceId]
			kpis := types.ParseKpis(alert.Kpis)

			var brokenKpis []types.KPI
			// shouldn't be necessary if alerts would only be sent when KPIs are violated
			for _, kpi := range kpis {
				if kpi.Level == "Broken" || kpi.Level == "Critical" {
					brokenKpis = append(brokenKpis, kpi)
				}
			}
			if len(brokenKpis) > 0 {
				rs.handleRoleDecisions(roles, brokenKpis, serviceId)
			}
		case <-ticker.C:
			//periodically check all services
			for serviceId, roles := range rolesByServiceId {
				kpis, err := rs.KpiRetriever.Get(serviceId)
				if err != nil {
					log.Printf("Failed to get KPIs for service %v: %v", serviceId, err)
					continue
				}
				rs.handleRoleDecisions(roles, kpis, serviceId)
			}
		}
	}
}

func (rs *RoleSelector) handleRoleDecisions(roles []*types.Role, kpis []types.KPI, serviceId string) {
	decisions, _ := rs.Policy.DecidePolicy(roles, kpis, []types.Resource{})
	for _, role := range roles {
		roleId := role.Id
		imageId := role.ImageId
		if shouldRun, ok := decisions[roleId]; ok {
			if shouldRun {
				// Only start the role if it's not running
				if !role.IsRunning {
					log.Printf("Decided to run role: %v, serviceId: %v", roleId, serviceId)
					rs.RoleRunner.Run(roleId, serviceId, imageId)
					role.IsRunning = true
				}
			} else {
				// Only stop the role if it's running
				if role.IsRunning {
					log.Printf("Decided to stop role: %v, serviceId: %v", roleId, serviceId)
					rs.RoleRunner.Stop(roleId, serviceId, imageId)
					role.IsRunning = false
				}
			}
		}
	}
}
