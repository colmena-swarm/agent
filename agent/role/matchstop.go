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
package role

import (
	"log"

	"colmena.bsc.es/agent/device"
)

func matchstop(device device.Info, roles []Role, kpiMatcher KpiMatcher) []Role {
	if (device.Strategy == "EAGER") {
		return []Role{} //eager devices never stop roles
	}

	matching := []Role{}
	for _, role := range roles {
		if (qosMet(role, kpiMatcher) && !hasWorkload(role, kpiMatcher, device)) {
			matching = append(matching, role)
		}
	}
	return matching
}

func qosMet(role Role, kpiMatcher KpiMatcher) bool {
	kpiMet := kpisMet(role.Kpis, kpiMatcher)
	log.Printf("Kpi met: roleId: %v, kpiMet: %v", role.Id, kpiMet)
	return kpiMet
}

func hasWorkload(role Role, kpiMatcher KpiMatcher, device device.Info) bool {
	kpi := KPI{
		Key:           	role.ServiceId + "/num_executions_" + device.Name + "_" + role.Id,
		Threshold:     	1,
		FromUnit:      	5,
		FromType:	  	"s",
		Comparison:		">=",
	}
	kpiMet := kpisMet([]KPI{kpi}, kpiMatcher)
	log.Printf("Has workload? roleId: %v, hasWorkload: %v", role.Id, kpiMet)
	return kpiMet
}
